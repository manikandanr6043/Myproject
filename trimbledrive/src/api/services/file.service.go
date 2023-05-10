package services

import (
	"time"

	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/model/fieldnames"
	"trimble.com/tdrive/api/model/file"
	"trimble.com/tdrive/api/utils"
)

// FileService -> struct for file service
type FileService struct {
	latestRepository *repository.LatestRepository
	fileSpaceService *FilespaceService
	folderService    *FolderService
	blobService      *BlobService
}

// NewFileService creates and returns an instance of FileService
func NewFileService(latestRepository *repository.LatestRepository, fileSpaceService *FilespaceService, folderService *FolderService, blobService *BlobService) *FileService {
	return &FileService{
		latestRepository: latestRepository,
		fileSpaceService: fileSpaceService,
		folderService:    folderService,
		blobService:      blobService,
	}
}

// GetFileDetails Fetches file details for a given fileId
// returns error if file not found or user doesn't have access or any issues.
func (f *FileService) GetFileDetails(ctx *requestcontext.RequestContext, spaceId string, fileId string, requestParams file.DetailsParams) (*file.DetailsResponse, *api_error.ApiError) {
	space, err := f.fileSpaceService.getAndValidateAccess(ctx, spaceId, nil)
	if err != nil {
		return nil, err
	}

	fileEntity, err := f.latestRepository.GetEntityWithParentFoldersList(ctx, spaceId, fileId, constants.File)
	if err != nil {
		return nil, err
	}
	latest := fileEntity.Latest
	if latest.Deleted && (requestParams.Deleted == nil || !*requestParams.Deleted) {
		return nil, api_error.ResourceNotFound
	}
	var topActiveParentFolderId *string
	if !latest.Deleted {
		topActiveParentFolderId = repository.GetTopActiveParentFolderId(fileEntity.ParentFolders, *space.RootId)
	}
	return f.convertToResponse(ctx, &latest, topActiveParentFolderId, requestParams.Fields, requestParams.UrlExpiryInMins)
}

// UpdateFile Updates a file
// Update can be either rename or move or undelete operation
// returns error if file not found/valid or user doesn't have access or any issues.
func (f *FileService) UpdateFile(ctx *requestcontext.RequestContext, spaceId string, fileId string, request file.UpdateFileRequest, headers file.HeaderParams, requestParams file.UpdateParams) (*file.DetailsResponse, *api_error.ApiError) {
	space, err := f.fileSpaceService.getAndValidateAccessForMutation(ctx, spaceId, false)
	if err != nil {
		return nil, err
	}

	// Name validation
	// This will not be performed when `enforceSafeNames` set explicitly to false
	if space.EnforceSafeNames {
		if validNameErr := utils.ValidName(request.Name); validNameErr != nil {
			return nil, validNameErr
		}
	}

	ifMatch, ifMatchErr := utils.ValidateAndGetIfMatchHeader(headers.IfMatch)
	if ifMatchErr != nil {
		return nil, ifMatchErr
	}

	fileEntity, err := f.latestRepository.GetEntityWithParentFoldersList(ctx, spaceId, fileId, constants.File)
	if err != nil {
		return nil, err
	}
	if fileEntity.Latest.Deleted && (requestParams.Deleted == nil || !*requestParams.Deleted) {
		return nil, api_error.ResourceNotFound
	}

	var topActiveParentFolderId *string
	if request.ParentFolderId != nil {
		destinationFolderEntity, destErr := f.latestRepository.GetEntityWithParentFoldersList(ctx, spaceId, *request.ParentFolderId, constants.Folder)
		if destErr != nil {
			return nil, destErr
		}
		if !destinationFolderEntity.Latest.Deleted {
			topActiveParentFolderId = repository.GetTopActiveParentFolderId(destinationFolderEntity.ParentFolders, *space.RootId)
		}
	} else {
		topActiveParentFolderId = repository.GetTopActiveParentFolderId(fileEntity.ParentFolders, *space.RootId)
	}

	updatedResult, updateErr := f.latestRepository.UpdateNameOrParentFolder(ctx, space, constants.File, fileId, request.Name, request.ParentFolderId, ifMatch)
	if updateErr != nil {
		return nil, updateErr
	}
	return f.convertToResponse(ctx, updatedResult, topActiveParentFolderId, requestParams.Fields, requestParams.UrlExpiryInMins)
}

// DeleteFile Deletes a file
// returns error if folder not found or user doesn't have access or any issues.
func (f *FileService) DeleteFile(ctx *requestcontext.RequestContext, spaceId string, fileId string, headers file.HeaderParams, requestParams file.RequestParams) (*file.DetailsResponse, *api_error.ApiError) {
	space, spaceErr := f.fileSpaceService.getAndValidateAccessForMutation(ctx, spaceId, false)
	if spaceErr != nil {
		return nil, spaceErr
	}

	ifMatch, ifMatchErr := utils.ValidateAndGetIfMatchHeader(headers.IfMatch)
	if ifMatchErr != nil {
		return nil, ifMatchErr
	}

	result, deleteErr := f.latestRepository.DeleteById(ctx, space, constants.File, fileId, ifMatch)
	if deleteErr != nil {
		return nil, deleteErr
	}
	return f.convertToResponse(ctx, result, nil, requestParams.Fields, nil)
}

func (f *FileService) convertToResponse(ctx *requestcontext.RequestContext, latestEntity *repository.Latest, topActiveParentFolderId *string, fields *[]string, urlExpiryInMins *int) (*file.DetailsResponse, *api_error.ApiError) {
	deleted := &latestEntity.Deleted
	if deleted != nil && !*deleted {
		deleted = nil
	}
	response := &file.DetailsResponse{
		Id:         *latestEntity.Id,
		Type:       latestEntity.Type,
		Version:    repository.FormatVersion(latestEntity.MajorVersion, latestEntity.MinorVersion),
		SpaceId:    latestEntity.SpaceId,
		CreatedBy:  latestEntity.CreatedBy,
		CreatedOn:  latestEntity.Created,
		ModifiedBy: latestEntity.ModifiedBy,
		ModifiedOn: latestEntity.Modified,
		Deleted:    deleted,
	}
	if fields != nil {
		for _, field := range *fields {
			switch field {
			case constants.DownloadUrl:
				downloadUrl, err := f.getFileDownloadUrl(ctx, latestEntity.Name, nil, latestEntity.Contents, urlExpiryInMins)
				if err != nil {
					return nil, err
				}
				response.DownloadUrl = downloadUrl
			case fieldnames.Name:
				response.Name = &latestEntity.Name
			case fieldnames.Size:
				response.Size = &latestEntity.Size
			case fieldnames.ParentFolderId:
				response.ParentFolderId = latestEntity.ParentFolderId
			case fieldnames.TopActiveParentFolderId:
				response.TopActiveParentFolderId = topActiveParentFolderId
			}
		}
	}

	return response, nil
}

func (f *FileService) getFileDownloadUrl(ctx *requestcontext.RequestContext, fileName string, format *string, contents []repository.Content, urlExpiryInMins *int) (*string, *api_error.ApiError) {
	var srcFileDownloadUrl *string
	var expiryTime time.Time
	if urlExpiryInMins == nil {
		expiryTime = time.Now().UTC().Add(time.Minute * constants.DefaultExpiryInMinutes)
	} else {
		expiryTime = time.Now().UTC().Add(time.Minute * time.Duration(*urlExpiryInMins))
	}
	var storagePath string
	for _, content := range contents {
		if content.Format == format {
			storagePath = content.Location
			break
		}
	}
	downloadUrl, err := f.blobService.GenerateDownloadUrl(ctx, storagePath, fileName, expiryTime)
	if err != nil {
		return nil, err
	}
	srcFileDownloadUrl = &downloadUrl
	return srcFileDownloadUrl, nil
}
