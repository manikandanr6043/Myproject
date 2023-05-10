package upload

import "time"

// FileUploadResponse defines model for FileUploadResponse.
type FileUploadResponse struct {

	// UploadId Upload identifier
	UploadId string `json:"uploadId"`

	// Status Overall status of the upload
	Status string `json:"status"`

	Input FileUploadInput `json:"input"`

	Result *FileUploadResult `json:"result,omitempty"`

	// Created Upload created time
	CreatedOn time.Time `json:"createdOn"`
	CreatedBy string    `json:"createdBy"`

	// Modified File saved time
	ModifiedOn time.Time `json:"modifiedOn"`

	// ErrorReason Error reason in case status is ERROR
	ErrorReason *string `json:"errorReason,omitempty"`
}

// FileUploadResult defines model for File Upload result details.
type FileUploadResult struct {
	// Version File version identifier assigned on upload completion. Property is returned only when status is DONE.
	Version *string `json:"version,omitempty"`
}

// FileUploadInput defines model for File Upload Input details.
type FileUploadInput struct {
	// SpaceId file space id of the file
	SpaceId string `json:"spaceId"`
	// Name file name
	Name *string `json:"name,omitempty"`
	// Contents content that has format not specified refers to source file
	Contents []ContentDetails `json:"contents"`

	// FileId File identifier
	FileId string `json:"fileId,omitempty"`

	// ParentFolderId Parent folder id of the file
	ParentFolderId string `json:"parentFolderId,omitempty"`
}

// ContentDetails content that has format not specified refers to source file
type ContentDetails struct {
	// Format Supported values are TRB, PDF and THUMBNAIL. For source file this property can be ignored
	Format *string `json:"format,omitempty"`

	// Status file upload status, the valid values are UPLOADABLE, UPLOADED.
	Status string `json:"status,omitempty"`

	// Type Upload mode, it maybe SINGLEPART or MULTIPART
	UploadMode string `json:"uploadMode,omitempty"`

	// Url File should be uploaded to this url with the PUT operation and no additional headers including Authorization to be sent. Returned only when status is UPLOADABLE
	Url *string `json:"url,omitempty"`

	// Size Content size
	Size *int64 `json:"size,omitempty"`

	// UpdatedOn Content Updated time
	UpdatedOn time.Time `json:"updatedOn"`
}
