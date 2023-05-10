package folder

// RequestParams defines common parameters for All Folder APIs
type RequestParams struct {

	// Fields Option to specify fields to be included in the response.
	Fields *[]string `form:"fields,omitempty" json:"fields,omitempty"`
}

// DeleteParams defines parameters for DeleteFolder.
type DeleteParams struct {
	// Recursive Presence of this parameter deletes the directory hierarchy, including files.
	Recursive *bool `form:"recursive,omitempty" json:"recursive,omitempty"`

	// RequestParams defines common request parameters for FOLDER APIs
	RequestParams
}

// DetailsParams defines parameters for GetFolderDetails.
type DetailsParams struct {
	Version *int64 `form:"version,omitempty" json:"version,omitempty"`

	// Deleted If specified and true, resource is returned even if it is marked as deleted. Otherwise 404 error returned for deleted resource.
	Deleted *bool `form:"deleted,omitempty" json:"deleted,omitempty"`

	// RequestParams defines common request parameters for FOLDER APIs
	RequestParams
}

// UpdateParams defines parameters for UpdateFolder.
type UpdateParams struct {
	// Deleted To be passed with value as `true` for restoring a deleted folder. Otherwise the request fails with 404 if resource is currently marked as deleted.
	Deleted *bool `form:"deleted,omitempty" json:"deleted,omitempty"`

	// RequestParams defines common request parameters for FOLDER APIs
	RequestParams
}

// ListParams defines parameters for ListFolderChildren.
type ListParams struct {

	// RequestParams defines common request parameters for FOLDER APIs
	RequestParams

	// Deleted If specified and true, children marked as deleted will be also included into response.
	Deleted *bool `form:"deleted" json:"deleted,omitempty"`

	// SkipToken Optimized cursor to fetch items. `skipToken` represents a position in the queried collection.
	// If not provided, first page of results will be returned.
	SkipToken *string `form:"skipToken,omitempty" json:"skipToken,omitempty"`

	// Top Maximum amount of items to return in the response.
	// It is not guaranteed that exact number of items will be returned, but up to the provided value as max.
	// The maximum value for this parameter is **`10000`** and the default value will be same as the maximum value.
	// If requested number of items exceeds the limit **10000**, then the request will fail with an error message.
	// Note that max items limit might be increased in the future, clients should not make any assumptions on the number of items returned by default.
	Top int64 `form:"top,default=10000" json:"top"`

	// OrderBy Sort order for returned items. Default order is by created time in descending order. Specify `createdOn asc` for ascending order.
	OrderBy *string `form:"orderBy" json:"orderBy"`
}

// ListVersionsParams defines parameters for ListFolderVersions.
type ListVersionsParams struct {
	// Fields Option to specify fields to be included in the response.
	Fields *[]string `form:"fields,omitempty" json:"fields,omitempty"`

	// Deleted If specified and true, children marked as deleted will be also included into response.
	Deleted *bool `form:"deleted" json:"deleted,omitempty"`

	FromVersion *string `form:"from" json:"from,omitempty"`

	ToVersion *string `form:"to" json:"to,omitempty"`

	// SkipToken Optimized cursor to fetch items. `skipToken` represents a position in the queried collection.
	// If not provided, first page of results will be returned.
	SkipToken *string `form:"skipToken,omitempty" json:"skipToken,omitempty"`

	// Top Maximum amount of items to return in the response.
	// It is not guaranteed that exact number of items will be returned, but up to the provided value as max.
	// The maximum value for this parameter is **`10000`** and the default value will be same as the maximum value.
	// If requested number of items exceeds the limit **10000**, then the request will fail with an error message.
	// Note that max items limit might be increased in the future, clients should not make any assumptions on the number of items returned by default.
	Top int64 `form:"top,default=10000" json:"top"`

	// OrderBy Sort order for returned items. Default order is by created time in descending order. Specify `createdOn asc` for ascending order.
	OrderBy *string `form:"orderBy" json:"orderBy"`
}

type HeaderParams struct {
	//IfMatch header param for version validation
	IfMatch *string `header:"If-Match,omitempty" json:"If-Match,omitempty"`
}
