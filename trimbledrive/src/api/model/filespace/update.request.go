package filespace

// UpdateRequest defines model for FileSpaceUpdateRequest.
type UpdateRequest struct {
	// Description The description of the file space.
	Description *string `json:"description,omitempty" binding:"omitempty,max=1024"`
}
