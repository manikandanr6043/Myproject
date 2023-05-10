package upload

// FileUploadParams defines parameters for InitiateFileUpload.
type FileUploadParams struct {
	// UrlExpiryInMins Number of minutes for which the pre-signed upload URL for the file will be valid.
	UrlExpiryInMins *int `form:"urlExpiryInMins,omitempty" json:"urlExpiryInMins,omitempty" binding:"omitempty,min=1,max=720"`
}

// FileUploadHeaders defines parameters for InitiateFileUpload.
type FileUploadHeaders struct {
	IfMatch *string `header:"If-Match,omitempty"`

	// IfNoneMatch If this header is passed with value `*` in the request the upload will fail if there is already file with same name exists in the same folder
	IfNoneMatch *string `header:"If-None-Match,omitempty"`
}

// GetFileUploadDetailsParams defines parameters for GetFileUploadDetails.
type GetFileUploadDetailsParams struct {
	// UrlExpiryInMins Number of minutes for which the pre-signed upload URL for the file will be valid.
	UrlExpiryInMins *int `form:"urlExpiryInMins,omitempty" json:"urlExpiryInMins,omitempty" binding:"omitempty,min=1,max=720"`
}
