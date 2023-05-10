package filespace

import "time"

// Response defines model for FileSpaceDetails.
type Response struct {

	// Id The globally unique identifier of the file space.
	Id *string `json:"id"`
	// Description The description of the file space.
	Description *string `json:"description,omitempty"`
	// ChangeToken The change token that represents a state of the space and all it's content (resources in the space) at specific moment of time when this resource has been returned.
	//Can be used for fetching the changes after this moment of time in incremental synchronisation scenarios by passing this token to `ListChanges` API.
	//The change in the value of this token means something has been changed in the space.
	ChangeToken *string `json:"changeToken"`
	// RootId Root folder identifier.
	RootId *string `json:"rootId"`
	// Deleted Indicates whether the space is in the deleted state. Deleted state makes repository content not accessible.
	Deleted *bool `json:"deleted,omitempty"`

	CreatedOn *time.Time `json:"createdOn"`
	CreatedBy *string    `json:"createdBy"`

	ModifiedOn *time.Time `json:"modifiedOn"`
	ModifiedBy *string    `json:"modifiedBy"`

	// EnforceSafeNames Option to manage file and folder naming restrictions.
	EnforceSafeNames bool `json:"enforceSafeNames"`
}
