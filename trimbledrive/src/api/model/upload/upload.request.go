package upload

// FileUploadRequest defines model for FileUploadRequest.
type FileUploadRequest struct {
	// Id The identifier for a new file can be given in the request. The id has to be unique in the file space and cannot be updated after creation. If id is not specified service will generate unique id.
	Id *string `json:"id,omitempty" binding:"omitempty,max=128,UrlSafe"`

	// Name file name
	Name *string `json:"name,omitempty" binding:"omitempty,max=255"`

	// ParentFolderId Parent folder identifier to which the file has to be uploaded
	ParentFolderId string `json:"parentFolderId" binding:"max=128,UrlSafe"`

	// Contents Content that has format not specified refers to source file. Object for source file in contents is not required for single-part uploads.
	Contents *[]FileUploadContent `json:"contents,omitempty"`
}

// FileUploadContent Content that has format not specified refers to source file. Object for source file in contents is not required for single-part uploads.
type FileUploadContent struct {
	// Format Supported values are TRB, PDF and THUMBNAIL. For source file this property can be ignored
	Format *string `json:"format,omitempty"`
}
