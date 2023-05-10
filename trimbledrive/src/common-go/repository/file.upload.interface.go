package repository

import (
	"time"

	"trimble.com/common/api_error"
	"trimble.com/common/requestcontext"
)

type FileUpload struct {
	UploadId    string        `json:"id"`
	Status      string        `json:"status"`
	Input       UploadInput   `json:"input"`
	Result      *UploadResult `json:"result"`
	CreatedOn   time.Time     `json:"createdOn"`
	CreatedBy   string        `json:"createdBy"`
	ModifiedOn  time.Time     `json:"modifiedOn,omitempty"`
	ErrorReason *string       `json:"errorReason,omitempty"`
}

type UploadInput struct {
	SpaceId        string                   `json:"spaceId"`
	Name           *string                  `json:"name"`
	FileId         string                   `json:"fileId"`
	ParentFolderId string                   `json:"parentFolderId"`
	Contents       map[string]UploadContent `json:"contents"`
	IfMatch        *Version                 `json:"ifMatch,omitempty"`
	IfNoneMatch    *string                  `json:"ifNoneMatch,omitempty"`
}

type UploadContent struct {
	Format     *string   `json:"format,omitempty"`
	Status     string    `json:"status"`
	UploadMode string    `json:"uploadMode"`
	UpdatedOn  time.Time `json:"updatedOn"`
	Size       *int64    `json:"size"`
}

type UploadResult struct {
	Version *string `json:"version,omitempty"`
}

type FileUploadRepositoryInterface interface {
	Insert(ctx *requestcontext.RequestContext, fileUpload *FileUpload) *api_error.ApiError
	FindByUploadId(ctx *requestcontext.RequestContext, uploadId string) (*FileUpload, *api_error.ApiError)
}
