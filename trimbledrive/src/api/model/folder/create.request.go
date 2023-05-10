package folder

// CreateRequest defines model for CreateFolderRequest.
type CreateRequest struct {
	// Id The identifier for a new folder can be given in the request. The id has to be unique in the file space and cannot be updated after creation. If id is not specified service will generate unique id.
	Id             *string `json:"id,omitempty" binding:"omitempty,max=128,UrlSafe"`
	Name           string  `json:"name" binding:"required,max=255"`
	ParentFolderId string  `json:"parentFolderId" binding:"required,max=128,UrlSafe"`
}
