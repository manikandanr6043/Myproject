package filespace

// Base params
type SpaceParams struct {

	// Deleted If specified and true, operation on resource is allowed even if it is marked as deleted. Otherwise, 404 error returned for deleted resource.
	Deleted *bool `form:"deleted,omitempty"`
}

// DetailsParams defines parameters for GetFilespaceDetails.
type DetailsParams struct {
	SpaceParams
}
