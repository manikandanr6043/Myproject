package folder

// UpdateRequest defines model for UpdateFolderRequest.
type UpdateRequest struct {
	// ParentFolderId The new parent folder identifier
	ParentFolderId *string `json:"parentFolderId,omitempty" binding:"omitempty,max=128,UrlSafe"`

	// Name The new name.
	Name *string `json:"name,omitempty" binding:"omitempty,max=255"`
}
