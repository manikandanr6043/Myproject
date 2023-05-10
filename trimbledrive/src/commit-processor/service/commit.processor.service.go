package service

import (
	"time"

	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
	"trimble.com/tdrive/commit-processor/model"
)

// CommitProcessorService -> struct for commit processor service
type CommitProcessorService struct {
	uploadRepository *repository.FileUploadRepository
	spaceRepository  *repository.FilespaceRepository
	latestRepository *repository.LatestRepository
}

// NewCommitProcessorService creates new CommitProcessorService instance
func NewCommitProcessorService(uploadRepository *repository.FileUploadRepository, spaceRepository *repository.FilespaceRepository, latestRepository *repository.LatestRepository) *CommitProcessorService {
	return &CommitProcessorService{
		uploadRepository: uploadRepository,
		spaceRepository:  spaceRepository,
		latestRepository: latestRepository,
	}
}

// CommitFileUpload Commits File Upload based on the event and returns false if the file commit fails and needs to be retries, true otherwise
func (c *CommitProcessorService) CommitFileUpload(rCtx *requestcontext.RequestContext, commitMessage model.CommitMessage) bool {
	uploadId := commitMessage.Data.UploadId
	// Get and Validate Upload is in UPLOADABLE state
	uploadEntity, err := c.uploadRepository.FindByUploadId(rCtx, uploadId)
	if err != nil || uploadEntity.Status != constants.FileUploadStatusUploadable {
		rCtx.Logger().Debug("Upload status not " + constants.FileUploadStatusUploadable)
		return true
	}
	rCtx.SetUserId(uploadEntity.CreatedBy)
	latestFileEntity, commitErr := c.validateAndCommit(rCtx, uploadEntity)
	if commitErr != nil {
		// Trigger retry on Internal Server Error
		if commitErr == api_error.InternalServerError {
			return false
		} else {
			updateErr := c.uploadRepository.UpdateUploadById(rCtx, uploadEntity, constants.FileUploadStatusError, time.Now().UTC(), nil, &commitErr.ErrorCode)
			if updateErr != nil {
				rCtx.Logger().Error("Error on updating File Upload status", zap.Error(updateErr))
				return true
			}
		}
	} else {
		version := repository.FormatVersion(latestFileEntity.MajorVersion, latestFileEntity.MinorVersion)
		updateErr := c.uploadRepository.UpdateUploadById(rCtx, uploadEntity, constants.FileUploadStatusDone, time.Now().UTC(), &version, nil)
		if updateErr != nil {
			rCtx.Logger().Error("Error on updating File Upload status", zap.Error(updateErr))
			return true
		}
	}
	return true
}

func (c *CommitProcessorService) getAndValidateAccessForMutation(ctx *requestcontext.RequestContext, spaceId string) (*repository.Filespace, *api_error.ApiError) {
	space, err := c.spaceRepository.IncrementClock(ctx, spaceId)
	if err != nil {
		return nil, err
	}

	if *space.Deleted {
		return nil, api_error.SpaceNotFound
	}

	return space, nil
}

func (c *CommitProcessorService) validateAndCommit(rCtx *requestcontext.RequestContext, uploadEntity *repository.FileUpload) (*repository.Latest, *api_error.ApiError) {
	spaceId := uploadEntity.Input.SpaceId
	space, err := c.getAndValidateAccessForMutation(rCtx, spaceId)
	if err != nil {
		return nil, err
	}
	// Skipping parent folder id check as its already validated during the POST call.
	// Parent active/deleted state does not affect the file upload to the parent
	// Get File By Id
	latestEntity, err := c.latestRepository.GetById(rCtx, spaceId, constants.File, uploadEntity.Input.FileId)
	if err != nil && err != api_error.ResourceNotFound {
		return nil, err
	}
	// Get SRC file size
	if latestEntity != nil {
		// Perform If-None-Match Validation
		ifNoneMatch := uploadEntity.Input.IfNoneMatch
		if ifNoneMatch != nil && *ifNoneMatch == "*" {
			return nil, api_error.FileExists
		}
		// Perform Update if existing file
		return c.latestRepository.UpdateFileContent(rCtx, space, uploadEntity.Input.FileId, uploadEntity.Input.ParentFolderId,
			getSrcFileSize(uploadEntity), convertToContentsEntity(uploadEntity), uploadEntity.Input.Name, uploadEntity.Input.IfMatch)
	}
	// Perform insert if no existing entry
	latest := &repository.Latest{
		Id:              &uploadEntity.Input.FileId,
		SpaceId:         spaceId,
		Name:            *uploadEntity.Input.Name,
		Type:            constants.File,
		ParentFolderId:  &uploadEntity.Input.ParentFolderId,
		Size:            getSrcFileSize(uploadEntity),
		Contents:        convertToContentsEntity(uploadEntity),
		CreatedOnClock:  space.Clock,
		ModifiedOnClock: space.Clock,
	}
	createErr := c.latestRepository.Create(rCtx, latest)
	return latest, createErr
}

func getSrcFileSize(uploadEntity *repository.FileUpload) int64 {
	var srcFileSize int64
	for _, v := range uploadEntity.Input.Contents {
		if v.Format == nil {
			srcFileSize = *v.Size
			break
		}
	}
	return srcFileSize
}

func convertToContentsEntity(uploadEntity *repository.FileUpload) []repository.Content {
	var contents []repository.Content
	for k, v := range uploadEntity.Input.Contents {
		content := repository.Content{
			Format:   v.Format,
			Location: k,
			Size:     *v.Size,
		}
		contents = append(contents, content)
	}
	return contents
}
