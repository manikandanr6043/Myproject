package services

import (
	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/api/utils"

	"trimble.com/tdrive/api/model/fieldnames"
	"trimble.com/tdrive/api/model/folder"
)

// FolderService -> struct for folder service
type FolderService struct {
	folderRepository *repository.LatestRepository
	fileSpaceService *FilespaceService
}

// NewFolderService creates new folder service instance
func NewFolderService(repository *repository.LatestRepository, fileSpaceService *FilespaceService) *FolderService {
	return &FolderService{
		folderRepository: repository,
		fileSpaceService: fileSpaceService,
	}
}

// CreateFolder Creates a folder under a given parent
// returns error if parent not found or user doesn't have access or any issues.
func (f *FolderService) CreateFolder(ctx *requestcontext.RequestContext, spaceId string, request folder.CreateRequest, requestParams folder.RequestParams) (*folder.Response, *api_error.ApiError) {
	space, spaceErr := f.fileSpaceService.getAndValidateAccessForMutation(ctx, spaceId, false)
	if spaceErr != nil {
		return nil, spaceErr
	}

	// Name validation
	// This will not be performed when `enforceSafeNames` set explicitly to false
	if space.EnforceSafeNames {
		if validNameErr := utils.ValidName(&request.Name); validNameErr != nil {
			return nil, validNameErr
		}
	}

	// Validate parent folder
	parentFolder, err := f.folderRepository.GetEntityWithParentFoldersList(ctx, spaceId, request.ParentFolderId, constants.Folder)
	if err != nil {
		return nil, err
	}
	if parentFolder.Latest.Deleted {
		return nil, api_error.ResourceNotFound
	}
	latest := &repository.Latest{
		Id:              request.Id,
		SpaceId:         spaceId,
		Name:            request.Name,
		Type:            constants.Folder,
		ParentFolderId:  &request.ParentFolderId,
		CreatedOnClock:  space.Clock,
		ModifiedOnClock: space.Clock,
	}
	err = f.folderRepository.Create(ctx, latest)
	if err != nil {
		return nil, err
	}

	topActiveParentFolderId := repository.GetTopActiveParentFolderId(parentFolder.ParentFolders, *space.RootId)
	return convertToResponse(latest, topActiveParentFolderId, requestParams.Fields), err
}

// UpdateFolder Updates a folder
// Update can be either rename or move or undelete operation
// returns error if folder not found/valid or user doesn't have access or any issues.
func (f *FolderService) UpdateFolder(ctx *requestcontext.RequestContext, spaceId string, id string, request folder.UpdateRequest, headers folder.HeaderParams, requestParams folder.UpdateParams) (*folder.Response, *api_error.ApiError) {
	space, err := f.fileSpaceService.getAndValidateAccessForMutation(ctx, spaceId, false)
	if err != nil {
		return nil, err
	}
	if id == *space.RootId {
		ctx.Logger().Debug("Not allowed to update root folder")
		return nil, api_error.InvalidFolderUpdateOperation
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

	var topActiveParentFolderId *string
	folderEntityResponse, err := f.folderRepository.GetEntityWithParentFoldersList(ctx, spaceId, id, constants.Folder)
	if err != nil {
		return nil, err
	}
	if folderEntityResponse.Latest.Deleted && (requestParams.Deleted == nil || !*requestParams.Deleted) {
		return nil, api_error.ResourceNotFound
	}
	if request.ParentFolderId != nil {
		destinationFolderResponse, err := f.folderRepository.GetEntityWithParentFoldersList(ctx, spaceId, *request.ParentFolderId, constants.Folder)
		if err != nil {
			return nil, err
		}
		if !destinationFolderResponse.Latest.Deleted {
			topActiveParentFolderId = repository.GetTopActiveParentFolderId(destinationFolderResponse.ParentFolders, *space.RootId)
		}
	} else {
		topActiveParentFolderId = repository.GetTopActiveParentFolderId(folderEntityResponse.ParentFolders, *space.RootId)
	}

	updatedResult, updateErr := f.folderRepository.UpdateNameOrParentFolder(ctx, space, constants.Folder, id, request.Name, request.ParentFolderId, ifMatch)
	if updateErr != nil {
		return nil, updateErr
	}
	return convertToResponse(updatedResult, topActiveParentFolderId, requestParams.Fields), err
}

// GetFolder Fetches folder details for a given folderId
// returns error if folder not found or user doesn't have access or any issues.
func (f *FolderService) GetFolder(ctx *requestcontext.RequestContext, spaceId string, folderId string, requestParams folder.DetailsParams) (*folder.Response, *api_error.ApiError) {
	space, err := f.fileSpaceService.getAndValidateAccess(ctx, spaceId, nil)
	if err != nil {
		return nil, err
	}

	folderEntityResponse, err := f.folderRepository.GetEntityWithParentFoldersList(ctx, spaceId, folderId, constants.Folder)
	if err != nil {
		return nil, err
	}
	if folderEntityResponse.Latest.Deleted && (requestParams.Deleted == nil || !*requestParams.Deleted) {
		return nil, api_error.ResourceNotFound
	}
	var topActiveParentFolderId *string
	if !folderEntityResponse.Latest.Deleted {
		topActiveParentFolderId = repository.GetTopActiveParentFolderId(folderEntityResponse.ParentFolders, *space.RootId)
	}
	return convertToResponse(&folderEntityResponse.Latest, topActiveParentFolderId, requestParams.Fields), nil
}

// DeleteFolder Deletes a folder
// returns error if folder not found or user doesn't have access or any issues.
func (f *FolderService) DeleteFolder(ctx *requestcontext.RequestContext, spaceId string, id string, headers folder.HeaderParams, requestParams folder.DeleteParams) (*folder.Response, *api_error.ApiError) {
	space, spaceErr := f.fileSpaceService.getAndValidateAccessForMutation(ctx, spaceId, false)
	if spaceErr != nil {
		return nil, spaceErr
	}
	if id == *space.RootId {
		ctx.Logger().Debug("Not allowed to delete root folder")
		return nil, api_error.InvalidFolderDeleteOperation
	}

	ifMatch, ifMatchErr := utils.ValidateAndGetIfMatchHeader(headers.IfMatch)
	if ifMatchErr != nil {
		return nil, ifMatchErr
	}

	result, deleteErr := f.folderRepository.DeleteById(ctx, space, constants.Folder, id, ifMatch)
	if deleteErr != nil {
		return nil, deleteErr
	}
	return convertToResponse(result, nil, requestParams.Fields), nil

}

// GetChildren Fetch all the children under given parent folder
// Returns list of children or error if any issues.
func (f *FolderService) GetChildren(ctx *requestcontext.RequestContext, spaceId string, folderId string, params folder.ListParams) (*folder.ListResponse, *api_error.ApiError) {
	// Validate file space
	space, spaceErr := f.fileSpaceService.getAndValidateAccess(ctx, spaceId, nil)
	if spaceErr != nil {
		return nil, spaceErr
	}

	// Fetch given folder
	givenFolderEntity, folderErr := f.folderRepository.GetEntityWithParentFoldersList(ctx, spaceId, folderId, constants.Folder)
	if folderErr != nil {
		return nil, folderErr
	}
	givenFolder := givenFolderEntity.Latest

	// Return 404 if given folder is deleted and deleted param is false
	if givenFolder.Deleted && (params.Deleted == nil || !*params.Deleted) {
		return nil, api_error.ResourceNotFound
	}

	// Validate max page size
	if err := utils.ValidateTopParameter(params.Top); err != nil {
		return nil, err
	}

	// Validate requested sort parameter
	sort, sortErr := utils.ValidateAndParseSortParameter(params.OrderBy, false)
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
	// Get immediate items under the given folder
	result, err := f.folderRepository.GetChildren(ctx, *space.ID, folderId, params.Deleted, cursor, *sort, params.Top)
	if err != nil {
		return nil, err
	}

	// Find out topActiveParentFolderId in case of active given folder
	var topActiveParentFolderId *string
	if !givenFolder.Deleted {
		topActiveParentFolderId = repository.GetTopActiveParentFolderId(givenFolderEntity.ParentFolders, *space.RootId)
	}

	// Convert result set to list response model
	response := convertToListResponse(result, topActiveParentFolderId, params.Fields)

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

// ValidateFolder Check if a given folder exists and active
// throws error if folder is not active or not found
func (f *FolderService) ValidateFolder(ctx *requestcontext.RequestContext, spaceId string, folderId string) (*repository.Latest, *api_error.ApiError) {
	// Validate folder
	folderEntity, folderErr := f.folderRepository.GetById(ctx, spaceId, constants.Folder, folderId)
	if folderErr != nil {
		return nil, folderErr
	}
	if folderEntity.Deleted {
		return nil, api_error.ResourceNotFound
	}
	return folderEntity, nil
}

// convertToListResponse Utility method to convert list of entities to list response
func convertToListResponse(items []*repository.Latest, topActiveParentFolderId *string, fields *[]string) *folder.ListResponse {

	// Empty results set (in case of empty results, response will be {"items":[]} )
	results := make([]folder.Response, 0)
	for _, item := range items {
		repository.ConvertToVisibleId(item)
		results = append(results, *convertToResponse(item, topActiveParentFolderId, fields))
	}
	return &folder.ListResponse{Items: results}

}

// convertToResponse Internal utility method to convert entity object to response
func convertToResponse(latest *repository.Latest, topActiveParentFolderId *string, fields *[]string) *folder.Response {
	deleted := &latest.Deleted
	if deleted != nil && !*deleted {
		deleted = nil
	}
	response := &folder.Response{
		Id:         *latest.Id,
		Version:    repository.FormatVersion(latest.MajorVersion, latest.MinorVersion),
		Type:       latest.Type,
		SpaceId:    latest.SpaceId,
		CreatedBy:  latest.CreatedBy,
		ModifiedBy: latest.ModifiedBy,
		CreatedOn:  latest.Created,
		ModifiedOn: latest.Modified,
		Deleted:    deleted,
	}
	if fields != nil {
		for _, field := range *fields {
			switch field {
			case fieldnames.Name:
				response.Name = &latest.Name
			case fieldnames.Size:
				response.Size = &latest.Size
			case fieldnames.ParentFolderId:
				response.ParentFolderId = latest.ParentFolderId
			case fieldnames.TopActiveParentFolderId:
				response.TopActiveParentFolderId = topActiveParentFolderId
			}
		}
	}

	return response

}
