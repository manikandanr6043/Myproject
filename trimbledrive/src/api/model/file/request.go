package file

// UpdateFileRequest defines request model for UpdateFileDetails API.
type UpdateFileRequest struct {

	// ParentFolderId The new parent folder identifier.
	ParentFolderId *string `json:"parentFolderId,omitempty" binding:"omitempty,max=128,UrlSafe"`

	// Name The new name.
	Name *string `json:"name,omitempty" binding:"omitempty,max=255"`
}
