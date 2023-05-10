package folder

import (
	"time"
)

// Response defines model for FolderDetailsResponse.
type Response struct {
	CreatedBy string    `json:"createdBy"`
	CreatedOn time.Time `json:"createdOn"`

	// ParentFolderId Folder identifier of the parent folder. Not present for root folder or for files not attached to any folder (attachments).
	ParentFolderId          *string   `json:"parentFolderId,omitempty"`
	Id                      string    `json:"id"`
	ModifiedBy              string    `json:"modifiedBy"`
	ModifiedOn              time.Time `json:"modifiedOn"`
	Name                    *string   `json:"name,omitempty"`
	Size                    *int64    `json:"size,omitempty"`
	SpaceId                 string    `json:"spaceId"`
	Type                    string    `json:"type"`
	Version                 string    `json:"version"`
	Deleted                 *bool     `json:"deleted,omitempty"`
	TopActiveParentFolderId *string   `json:"topActiveParentFolderId,omitempty"`
}

// ListResponse defines model for ListFolderDetailsResponse.
type ListResponse struct {
	// Items List of folder details
	Items []Response `json:"items"`

	// Next Responses that include only a partial set of the items identified by the request will contain a link that allows retrieving the next partial set of items. <br/>This link is called a next link. Clients should treat the URL of the next link as opaque, and should not append system query options to the URL of a next link. <br/>The url returned by server might be an absolute url or path with query parameters. <br/>If it is a path, client is responsible for adding a hostname to form a full url. <br/>Otherwise no transformations must be applied by client to the url and it should be sent as is to the server to fetch a next page of events.
	Next *string `json:"next,omitempty"`
}
