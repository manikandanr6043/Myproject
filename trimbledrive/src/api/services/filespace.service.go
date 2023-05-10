package services

import (
	"strconv"
	"time"

	"github.com/google/uuid"

	"trimble.com/common/api_error"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"

	"trimble.com/tdrive/api/model/filespace"
)

// FilespaceService -> struct for filespace service
type FilespaceService struct {
	repository *repository.FilespaceRepository
}

// NewFilespaceService creates new FilespaceService
func NewFilespaceService(repository *repository.FilespaceRepository) *FilespaceService {
	return &FilespaceService{
		repository: repository,
	}
}

// CreateFilespace Performs creation of file space
// returns errors if any issues in creation of filespace
func (f *FilespaceService) CreateFilespace(ctx *requestcontext.RequestContext, createRequestBody filespace.CreateRequest) (*filespace.Response, *api_error.ApiError) {
	newSpaceEntity := convertToFilespaceEntity(createRequestBody, ctx.UserId())
	// Save in database collection
	err := f.repository.Insert(ctx, newSpaceEntity)
	if err != nil {
		return nil, err
	}

	return convertToFilespaceResponse(newSpaceEntity), nil

}

// GetFilespace Performs Get filespace details
// returns errors if any issues in fetching filespace
func (f *FilespaceService) GetFilespace(ctx *requestcontext.RequestContext, spaceId string, requestParams filespace.DetailsParams) (*filespace.Response, *api_error.ApiError) {
	space, err := f.getAndValidateAccess(ctx, spaceId, requestParams.Deleted)
	if err != nil {
		return nil, err
	}

	return convertToFilespaceResponse(space), nil
}

// UpdateFilespace Performs Update filespace details
// returns errors if any issues in updating filespace
func (f *FilespaceService) UpdateFilespace(ctx *requestcontext.RequestContext, spaceId string, updateRequest filespace.UpdateRequest, updateParams filespace.SpaceParams) (*filespace.Response, *api_error.ApiError) {
	space, err := f.getAndValidateAccess(ctx, spaceId, updateParams.Deleted)
	if err != nil {
		return nil, err
	}

	// Handle empty request
	if updateRequest.Description == nil {
		return convertToFilespaceResponse(space), nil
	}

	// Updated entry in database
	space.Description = updateRequest.Description
	modified := time.Now().UTC()
	userId := ctx.UserId()
	space.ModifiedBy = &userId
	space.Modified = &modified
	space, updateErr := f.repository.UpdateDescription(ctx, space)
	if updateErr != nil {
		return nil, updateErr
	}

	return convertToFilespaceResponse(space), nil
}

// DeleteFilespace Performs delete filespace
// returns errors if any issues in deleting filespace
func (f *FilespaceService) DeleteFilespace(ctx *requestcontext.RequestContext, spaceId string) *api_error.ApiError {
	space, err := f.getAndValidateAccess(ctx, spaceId, nil)
	if err != nil {
		return err
	}

	// Updated entry in database
	deleted := true
	modified := time.Now().UTC()
	userId := ctx.UserId()
	space.Deleted = &deleted
	space.ModifiedBy = &userId
	space.Modified = &modified
	_, updateErr := f.repository.UpdateDeleted(ctx, space)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

// getAndValidateAccess validate if user has access to the requested space
// returns repository.Filespace if validation is successful else api_error.ApiError is returned
func (f *FilespaceService) getAndValidateAccess(ctx *requestcontext.RequestContext, spaceId string, ignoreDeletedState *bool) (*repository.Filespace, *api_error.ApiError) {
	space, getErr := f.repository.GetById(ctx, spaceId)
	if getErr != nil {
		return nil, getErr
	}
	// Validate if user has access to space
	// TODO: Enhance validations further with ACL as part of: TD-58
	if *space.CreatedBy != ctx.UserId() {
		return nil, api_error.PermissionDenied
	}

	if (ignoreDeletedState == nil || !*ignoreDeletedState) && *space.Deleted {
		return nil, api_error.SpaceNotFound
	}

	return space, nil
}

func (f *FilespaceService) getAndValidateAccessForMutation(ctx *requestcontext.RequestContext, spaceId string, ignoreDeletedState bool) (*repository.Filespace, *api_error.ApiError) {
	space, err := f.repository.IncrementClock(ctx, spaceId)
	if err != nil {
		return nil, err
	}

	// Validate if user has access to space
	// TODO: Enhance validations further with ACL as part of: TD-58
	if *space.CreatedBy != ctx.UserId() {
		return nil, api_error.PermissionDenied
	}

	if !ignoreDeletedState && *space.Deleted {
		return nil, api_error.SpaceNotFound
	}

	return space, nil
}

// convertToFilespaceEntity Convert the FileSpaceRequest to Filespace entity
func convertToFilespaceEntity(createRequestBody filespace.CreateRequest, userId string) *repository.Filespace {
	spaceEntity := repository.NewFilespace(createRequestBody.Description, time.Now().UTC(), userId)
	// Generate space Id if not provided as part of the request
	if createRequestBody.Id == nil {
		id := uuid.NewString()
		spaceEntity.ID = &id
	} else {
		spaceEntity.ID = createRequestBody.Id
	}
	if createRequestBody.EnforceSafeNames != nil {
		spaceEntity.EnforceSafeNames = *createRequestBody.EnforceSafeNames
	} else {
		spaceEntity.EnforceSafeNames = true
	}
	return spaceEntity
}

// convertToFilespaceResponse Convert the Filespace Entity to Filespace Response
func convertToFilespaceResponse(spaceEntity *repository.Filespace) *filespace.Response {
	changeToken := strconv.FormatInt(*spaceEntity.Clock, 10)
	deleted := spaceEntity.Deleted
	if deleted != nil && !*deleted {
		deleted = nil
	}
	return &filespace.Response{
		Id:               spaceEntity.ID,
		Description:      spaceEntity.Description,
		ChangeToken:      &changeToken,
		RootId:           spaceEntity.RootId,
		Deleted:          deleted,
		CreatedOn:        spaceEntity.Created,
		CreatedBy:        spaceEntity.CreatedBy,
		ModifiedOn:       spaceEntity.Modified,
		ModifiedBy:       spaceEntity.ModifiedBy,
		EnforceSafeNames: spaceEntity.EnforceSafeNames,
	}
}
