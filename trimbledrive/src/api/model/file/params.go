package file

// RequestParams defines common parameters for All FILE APIs
type RequestParams struct {

	// Fields Option to specify fields to be included in the response.
	Fields *[]string `form:"fields,omitempty" json:"fields,omitempty"`
}

// DetailsParams defines parameters for GetFileDetails.
type DetailsParams struct {
	// UrlExpiryInMins Number of minutes for which the pre-signed Download URL for the file will be valid.
	UrlExpiryInMins *int `form:"urlExpiryInMins,omitempty" json:"urlExpiryInMins,omitempty" binding:"omitempty,min=1,max=720"`

	// Deleted If specified and true, resource is returned even if it is marked as deleted. Other wise 404 error returned for deleted resource.
	Deleted *bool `form:"deleted,omitempty" json:"deleted,omitempty"`

	// RequestParams defines common request parameters for FILE APIs
	RequestParams
}

// UpdateParams defines parameters for UpdateFileDetails.
type UpdateParams struct {
	// Deleted To be passed with value as `true` for restoring a deleted file.
	Deleted *bool `form:"deleted,omitempty" json:"deleted,omitempty"`

	// RequestParams defines common request parameters for FILE APIs
	RequestParams

	// UrlExpiryInMins Number of minutes for which the pre-signed Download URL for the file will be valid.
	UrlExpiryInMins *int `form:"urlExpiryInMins,omitempty" json:"urlExpiryInMins,omitempty" binding:"omitempty,min=1,max=720"`
}

type ListParams struct {
	// RequestParams defines common request parameters for FILE APIs
	RequestParams

	// Deleted If specified and true, deleted resource state is ignored.
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

	// UrlExpiryInMins Number of minutes for which the pre-signed Download URL for the file will be valid.
	UrlExpiryInMins *int `form:"urlExpiryInMins,omitempty" json:"urlExpiryInMins,omitempty" binding:"omitempty,min=1,max=720"`
}

// ListVersionsParams defines parameters for ListFileVersions.
type ListVersionsParams struct {
	// RequestParams defines common request parameters for FILE APIs
	RequestParams

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

	// UrlExpiryInMins Number of minutes for which the pre-signed Download URL for the file will be valid.
	UrlExpiryInMins *int `form:"urlExpiryInMins,omitempty" json:"urlExpiryInMins,omitempty" binding:"omitempty,min=1,max=720"`
}

// DeleteParams defines parameters for DeleteFile.
type DeleteParams struct {
	// Fields Option to specify fields to be included in the response.
	Fields *[]string `form:"fields,omitempty" json:"fields,omitempty"`
}

type HeaderParams struct {
	//IfMatch header param for version validation
	IfMatch *string `header:"If-Match,omitempty" json:"If-Match,omitempty"`
}
