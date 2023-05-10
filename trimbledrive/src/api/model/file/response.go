package file

import "time"

// DetailsResponse defines model for FileDetailsResponse.
type DetailsResponse struct {
	Id                      string    `json:"id"`
	Type                    string    `json:"type"`
	Version                 string    `json:"version"`
	SpaceId                 string    `json:"spaceId"`
	ParentFolderId          *string   `json:"parentFolderId,omitempty"`
	Name                    *string   `json:"name,omitempty"`
	Size                    *int64    `json:"size,omitempty"`
	CreatedBy               string    `json:"createdBy"`
	CreatedOn               time.Time `json:"createdOn"`
	ModifiedBy              string    `json:"modifiedBy"`
	ModifiedOn              time.Time `json:"modifiedOn"`
	ThumbnailUrl            *string   `json:"thumbnailUrl,omitempty"`
	DownloadUrl             *string   `json:"downloadUrl,omitempty"`
	Deleted                 *bool     `json:"deleted,omitempty"`
	TopActiveParentFolderId *string   `json:"topActiveParentFolderId,omitempty"`
}

// ListResponse defines model for ListFileDetailsResponse.
type ListResponse struct {
	// Items List of folder details
	Items []DetailsResponse `json:"items"`

	// Next Responses that include only a partial set of the items identified by the request will contain a link that allows retrieving the next partial set of items. <br/>This link is called a next link. Clients should treat the URL of the next link as opaque, and should not append system query options to the URL of a next link. <br/>The url returned by server might be an absolute url or path with query parameters. <br/>If it is a path, client is responsible for adding a hostname to form a full url. <br/>Otherwise no transformations must be applied by client to the url and it should be sent as is to the server to fetch a next page of events.
	Next *string `json:"next,omitempty"`
}
