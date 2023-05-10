package services

import (
	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/model/fieldnames"
	"trimble.com/tdrive/api/model/file"
	"trimble.com/tdrive/api/utils"
)

// FileVersionService -> struct for file versions service
type FileVersionService struct {
	resourceVersionRepository *repository.ResourceVersionRepository
	fileRepository            *repository.LatestRepository
	fileSpaceService          *FilespaceService
	fileService               *FileService
}

// NewFileVersionService creates and returns an instance of FileVersionService
func NewFileVersionService(resourceVersionRepository *repository.ResourceVersionRepository, fileRepository *repository.LatestRepository,
	fileSpaceService *FilespaceService, fileService *FileService) *FileVersionService {
	return &FileVersionService{
		resourceVersionRepository: resourceVersionRepository,
		fileRepository:            fileRepository,
		fileSpaceService:          fileSpaceService,
		fileService:               fileService,
	}
}

// GetFileVersionDetails Fetches file details for a given fileId and version
// returns error if file not found or user doesn't have access or any issues.
func (f *FileVersionService) GetFileVersionDetails(ctx *requestcontext.RequestContext, spaceId string, fileId string, version repository.Version, requestParams file.DetailsParams) (*file.DetailsResponse, *api_error.ApiError) {
	_, err := f.fileSpaceService.getAndValidateAccess(ctx, spaceId, nil)
	if err != nil {
		return nil, err
	}

	if versions, versionErr := f.resourceVersionRepository.GetResourceVersionById(ctx, spaceId, fileId, version, constants.File); versionErr != nil {
		return nil, versionErr
	} else {
		return f.convertVersionsToResponse(ctx, versions, requestParams.Fields, requestParams.UrlExpiryInMins)
	}
}

// ListFileVersions Fetch all the versions of given file
// Returns list of versions or error if any issues.
func (f *FileVersionService) ListFileVersions(ctx *requestcontext.RequestContext, spaceId string, fileId string, params file.ListVersionsParams) (*file.ListResponse, *api_error.ApiError) {
	// Validate file space
	space, spaceErr := f.fileSpaceService.getAndValidateAccess(ctx, spaceId, nil)
	if spaceErr != nil {
		return nil, spaceErr
	}

	// Validate max page size
	if err := utils.ValidateTopParameter(params.Top); err != nil {
		return nil, err
	}

	// Parse requested sort parameter
	sort, sortErr := utils.ValidateAndParseSortParameter(params.OrderBy, true)
	if sortErr != nil {
		return nil, sortErr
	}

	// Validate and parse skipToken to locate the cursor if present.
	var cursor *repository.Cursor
	if params.SkipToken != nil {
		if value, cursorErr := utils.ValidateAndParseSkipToken(ctx, spaceId, *params.SkipToken); cursorErr != nil {
			return nil, cursorErr
		} else {
			cursor = value
		}
	}

	// Fetch given file
	givenFile, err := f.fileRepository.GetById(ctx, spaceId, constants.File, fileId)
	if err != nil {
		return nil, err
	}
	// Return 404 if given file is deleted and deleted param is false
	if givenFile.Deleted && (params.Deleted == nil || !*params.Deleted) {
		return nil, api_error.ResourceNotFound
	}

	fromVersion, toVersion, versionErr := getFromAndToVersion(params.FromVersion, params.ToVersion)
	if versionErr != nil {
		return nil, versionErr
	}
	// Get list of versions for the given fileId
	result, err := f.resourceVersionRepository.ListResourceVersions(ctx, *space.ID, fileId, constants.File, fromVersion,
		toVersion, cursor, *sort, params.Top)
	if err != nil {
		return nil, err
	}

	// Convert result set to list response model
	items := make([]file.DetailsResponse, 0)
	for _, item := range result {
		repository.ConvertResourceVersionToVisibleId(item)
		res, err := f.convertVersionsToResponse(ctx, item, params.Fields, params.UrlExpiryInMins)
		if err != nil {
			return nil, err
		}
		items = append(items, *res)
	}

	response := &file.ListResponse{Items: items}

	resultSize := int64(len(result))
	if resultSize < params.Top {
		return response, nil
	}

	// Generate next url to access next page with relevant parameters
	lastItem := result[resultSize-1]
	next, nextErr := utils.GenerateNextUrl(ctx, spaceId, lastItem, sort)
	if nextErr != nil {
		return nil, nextErr
	}
	response.Next = next
	return response, nil
}

func (f *FileVersionService) convertVersionsToResponse(ctx *requestcontext.RequestContext, versionEntity *repository.ResourceVersion, fields *[]string, urlExpiryInMins *int) (*file.DetailsResponse, *api_error.ApiError) {
	deleted := &versionEntity.Deleted
	if deleted != nil && !*deleted {
		deleted = nil
	}
	response := &file.DetailsResponse{
		Id:         versionEntity.Id,
		Type:       versionEntity.Type,
		Version:    repository.FormatVersion(versionEntity.MajorVersion, versionEntity.MinorVersion),
		SpaceId:    versionEntity.SpaceId,
		CreatedBy:  versionEntity.CreatedBy,
		CreatedOn:  versionEntity.Created,
		ModifiedBy: versionEntity.ModifiedBy,
		ModifiedOn: versionEntity.Modified,
		Deleted:    deleted,
	}
	if fields != nil {
		for _, field := range *fields {
			switch field {
			case fieldnames.Name:
				response.Name = &versionEntity.Name
			case fieldnames.Size:
				response.Size = &versionEntity.Size
			case fieldnames.ParentFolderId:
				response.ParentFolderId = &versionEntity.ParentFolderId
			case constants.DownloadUrl:
				downloadUrl, err := f.fileService.getFileDownloadUrl(ctx, versionEntity.Name, nil, versionEntity.Contents, urlExpiryInMins)
				if err != nil {
					return nil, err
				}
				response.DownloadUrl = downloadUrl

			}
		}
	}

	return response, nil
}

func getFromAndToVersion(fromParam *string, toParam *string) (*repository.Version, *repository.Version, *api_error.ApiError) {
	fromVersion, versionErr := utils.ValidateAndGetVersion(fromParam)
	if versionErr != nil {
		return nil, nil, versionErr
	}
	toVersion, versionErr := utils.ValidateAndGetVersion(toParam)
	if versionErr != nil {
		return nil, nil, versionErr
	}
	return fromVersion, toVersion, nil
}
