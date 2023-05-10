package repository

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/data/azcosmos"
	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/requestcontext"
)

// FileUploadRepository struct for file upload repository
type FileUploadRepository struct {
	containerClient *azcosmos.ContainerClient
}

// NewFileUploadRepository creates instance of FileUploadRepository
func NewFileUploadRepository(cosmosDatabaseClient *azcosmos.DatabaseClient) *FileUploadRepository {
	containerClient, err := cosmosDatabaseClient.NewContainer("file_upload")
	if err != nil {
		log.Fatalf("Error creating cosmos container client: %s", err)
	}
	return &FileUploadRepository{
		containerClient: containerClient,
	}
}

// Insert create item in the FileUpload collection
func (f FileUploadRepository) Insert(ctx *requestcontext.RequestContext, fileUpload *FileUpload) *api_error.ApiError {
	pk := azcosmos.NewPartitionKeyString(fileUpload.UploadId)
	jsonEncodedEntity, encodeErr := jsonEncodeFileUploadEntity(ctx, fileUpload)
	if encodeErr != nil {
		return encodeErr
	}
	// Create item
	_, err := f.containerClient.CreateItem(ctx.DbCtx(), pk, jsonEncodedEntity, nil)
	if err != nil {
		ctx.Logger().Error("Error on CreateItem ", zap.Error(err))
		return api_error.InternalServerError
	}
	ctx.Logger().Debug("Inserted item", zap.Any("uploadID", fileUpload.UploadId))
	return nil
}

// FindByUploadId returns the FileUpload entity by the given uploadId is exists else returns error
func (f FileUploadRepository) FindByUploadId(ctx *requestcontext.RequestContext, uploadId string) (*FileUpload, *api_error.ApiError) {
	pk := azcosmos.NewPartitionKeyString(uploadId)
	// Read an item
	itemResponse, err := f.containerClient.ReadItem(ctx.DbCtx(), pk, uploadId, nil)
	if err != nil {
		ctx.Logger().Debug("File Upload not found", zap.String("uploadId", uploadId), zap.Error(err))
		return nil, api_error.UploadNotFound
	}
	return jsonUnmarshalFileUploadEntity(ctx, itemResponse.Value)
}

// UpdateUploadById updates the FileUpload entity with given arguments
func (f FileUploadRepository) UpdateUploadById(ctx *requestcontext.RequestContext, fileUpload *FileUpload, status string, modifiedOn time.Time, version *string, errorReason *string) *api_error.ApiError {
	pk := azcosmos.NewPartitionKeyString(fileUpload.UploadId)

	fileUpload.Status = status
	fileUpload.ModifiedOn = modifiedOn
	if version != nil {
		result := &UploadResult{
			Version: version,
		}
		fileUpload.Result = result
	}
	if errorReason != nil {
		fileUpload.ErrorReason = errorReason
	}
	jsonEncodedEntity, encodeErr := jsonEncodeFileUploadEntity(ctx, fileUpload)
	if encodeErr != nil {
		return encodeErr
	}
	// Patch an item
	_, err := f.containerClient.ReplaceItem(ctx.DbCtx(), pk, fileUpload.UploadId, jsonEncodedEntity, nil)
	if err != nil {
		ctx.Logger().Debug("Error on File Upload Patch", zap.String("uploadId", fileUpload.UploadId), zap.Error(err))
		return api_error.InternalServerError
	}
	return nil
}

// jsonEncodeFileUploadEntity returns json encoded form of the given FileUpload struct
func jsonEncodeFileUploadEntity(ctx *requestcontext.RequestContext, fileUpload *FileUpload) ([]byte, *api_error.ApiError) {
	jsonEncoded, err := json.Marshal(&fileUpload)
	if err != nil {
		ctx.Logger().Error("Error on jsonEncodeFileUploadEntity ", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	return jsonEncoded, nil
}

// jsonUnmarshalFileUploadEntity returns the FileUpload struct from the given json encoded form
func jsonUnmarshalFileUploadEntity(ctx *requestcontext.RequestContext, fileUploadEncoded []byte) (*FileUpload, *api_error.ApiError) {
	var fileUpload FileUpload
	err := json.Unmarshal(fileUploadEncoded, &fileUpload)
	if err != nil {
		ctx.Logger().Error("Error on jsonUnmarshallFileUploadEntity ", zap.Error(err))
		return nil, api_error.InternalServerError
	}
	return &fileUpload, nil
}
