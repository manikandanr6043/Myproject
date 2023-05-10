package services

import (
	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/model/fieldnames"
	"trimble.com/tdrive/api/model/folder"
	"trimble.com/tdrive/api/utils"
)

// FolderVersionService -> struct for folder version service
type FolderVersionService struct {
	resourceVersionRepository *repository.ResourceVersionRepository
	folderRepository          *repository.LatestRepository
	fileSpaceService          *FilespaceService
}

// NewFolderVersionService creates new folder version service instance
func NewFolderVersionService(resourceVersionRepository *repository.ResourceVersionRepository, folderRepository *repository.LatestRepository,
	fileSpaceService *FilespaceService) *FolderVersionService {
	return &FolderVersionService{
		resourceVersionRepository: resourceVersionRepository,
		folderRepository:          folderRepository,
		fileSpaceService:          fileSpaceService,
	}
}

// GetFolderVersionDetails Fetches folder details for a given folderId and version
// returns error if file not found or user doesn't have access or any issues.
func (f *FolderVersionService) GetFolderVersionDetails(ctx *requestcontext.RequestContext, spaceId string, folderId string, version repository.Version, params folder.DetailsParams) (*folder.Response, *api_error.ApiError) {
	_, err := f.fileSpaceService.getAndValidateAccess(ctx, spaceId, nil)
	if err != nil {
		return nil, err
	}

	if versions, versionErr := f.resourceVersionRepository.GetResourceVersionById(ctx, spaceId, folderId, version, constants.Folder); versionErr != nil {
		return nil, versionErr
	} else {
		return convertVersionsEntityToResponse(versions, params.Fields), nil
	}

}

// ListFolderVersions Fetch all the versions of given folder
// Returns list of versions or error if any issues.
func (f *FolderVersionService) ListFolderVersions(ctx *requestcontext.RequestContext, spaceId string, folderId string, params folder.ListVersionsParams) (*folder.ListResponse, *api_error.ApiError) {
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

	// Fetch given folder
	givenFolder, folderErr := f.folderRepository.GetById(ctx, spaceId, constants.Folder, folderId)
	if folderErr != nil {
		return nil, folderErr
	}
	// Return 404 if given folder is deleted and deleted param is false
	if givenFolder.Deleted && (params.Deleted == nil || !*params.Deleted) {
		return nil, api_error.ResourceNotFound
	}

	fromVersion, toVersion, versionErr := getFromAndToVersion(params.FromVersion, params.ToVersion)
	if versionErr != nil {
		return nil, versionErr
	}

	// Get list of versions for the given folderId
	result, err := f.resourceVersionRepository.ListResourceVersions(ctx, *space.ID, folderId, constants.Folder,
		fromVersion, toVersion, cursor, *sort, params.Top)
	if err != nil {
		return nil, err
	}

	// Convert result set to list response model
	items := make([]folder.Response, 0)
	for _, item := range result {
		repository.ConvertResourceVersionToVisibleId(item)
		items = append(items, *convertVersionsEntityToResponse(item, params.Fields))
	}

	response := &folder.ListResponse{Items: items}

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

// convertVersionsEntityToResponse Internal utility method to convert versions entity object to response
func convertVersionsEntityToResponse(versions *repository.ResourceVersion, fields *[]string) *folder.Response {
	deleted := &versions.Deleted
	if deleted != nil && !*deleted {
		deleted = nil
	}
	response := &folder.Response{
		Id:         versions.Id,
		Version:    repository.FormatVersion(versions.MajorVersion, versions.MinorVersion),
		Type:       versions.Type,
		SpaceId:    versions.SpaceId,
		CreatedBy:  versions.CreatedBy,
		ModifiedBy: versions.ModifiedBy,
		CreatedOn:  versions.Created,
		ModifiedOn: versions.Modified,
		Deleted:    deleted,
	}
	if fields != nil {
		for _, field := range *fields {
			switch field {
			case fieldnames.Name:
				response.Name = &versions.Name
			case fieldnames.Size:
				response.Size = &versions.Size
			case fieldnames.ParentFolderId:
				response.ParentFolderId = &versions.ParentFolderId
			}
		}
	}
	return response

}
