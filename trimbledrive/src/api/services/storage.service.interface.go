package services

import (
	"time"

	"trimble.com/common/api_error"
	"trimble.com/common/requestcontext"
)

type StorageServiceInterface interface {
	GenerateUploadUrl(ctx *requestcontext.RequestContext, storagePath string, fileName string, expiryTime time.Time) (string, *api_error.ApiError)
}
