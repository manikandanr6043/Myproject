package filespace

// CreateRequest defines model for CreateFileSpaceRequest.
type CreateRequest struct {

	// Id The identifier for a new file space can be given in the request.
	//The id has to be unique and cannot be updated after creation.
	//If id is not specified service will generate unique id.
	Id *string `json:"id,omitempty" binding:"omitempty,max=128,UrlSafe"`

	// Description The description of the file space.
	Description *string `json:"description,omitempty" binding:"omitempty,max=1024"`

	// EnforceSafeNames Option to manage file and folder naming restrictions.
	EnforceSafeNames *bool `json:"enforceSafeNames"`
}
