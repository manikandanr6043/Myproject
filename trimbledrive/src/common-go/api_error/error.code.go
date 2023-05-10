package api_error

import (
	"net/http"
	"strconv"

	"trimble.com/common/constants"
)

var (
	InternalServerError          = &ApiError{ErrorMessage: "Internal Server Error", StatusCode: http.StatusInternalServerError, ErrorCode: "InternalServerError"}
	UnAuthorized                 = &ApiError{ErrorMessage: "UnAuthorized", StatusCode: http.StatusUnauthorized, ErrorCode: "UnAuthorized"}
	SpaceNotFound                = &ApiError{ErrorMessage: "filespace not found", StatusCode: http.StatusNotFound, ErrorCode: "SpaceNotFound"}
	PermissionDenied             = &ApiError{ErrorMessage: "Permission denied for the requested resource", StatusCode: http.StatusForbidden, ErrorCode: "PermissionDenied"}
	DuplicateKey                 = &ApiError{ErrorMessage: "item with given id or name already exists", StatusCode: http.StatusConflict, ErrorCode: "DuplicateKey"}
	ResourceNotFound             = &ApiError{ErrorMessage: "file or folder not found", StatusCode: http.StatusNotFound, ErrorCode: "ResourceNotFound"}
	VersionNotFound              = &ApiError{ErrorMessage: "resource version not found", StatusCode: http.StatusNotFound, ErrorCode: "VersionNotFound"}
	UploadNotFound               = &ApiError{ErrorMessage: "upload not found", StatusCode: http.StatusNotFound, ErrorCode: "UploadNotFound"}
	InvalidFolderDeleteOperation = &ApiError{ErrorMessage: "folder delete operation cannot be performed on root folder", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidFolderDeleteOperation"}
	DuplicateResourceName        = &ApiError{ErrorMessage: "file or folder with same name already exists", StatusCode: http.StatusConflict, ErrorCode: "DuplicateResourceName"}
	DuplicateFormat              = &ApiError{ErrorMessage: "Duplicate values for format is not allowed", StatusCode: http.StatusBadRequest, ErrorCode: "DuplicateFormat"}
	InvalidFolderUpdateOperation = &ApiError{ErrorMessage: "folder update operation cannot be performed on root folder", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidFolderUpdateOperation"}
	InvalidIfMatchHeader         = &ApiError{ErrorMessage: "Invalid If-Match Header", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidIfMatchHeader"}
	InvalidUploadPayload         = &ApiError{ErrorMessage: "One of id or name is mandatory.In case of new file upload name is mandatory", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidPayload"}
	InvalidVersion               = &ApiError{ErrorMessage: "Specified version is not the latest version", StatusCode: http.StatusPreconditionFailed, ErrorCode: "InvalidVersion"}
	InvalidVersionParam          = &ApiError{ErrorMessage: "Please provide a valid version", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidParam"}
	FileWithSameNameExists       = &ApiError{ErrorMessage: "File with same name already exists under the given parent folder", ErrorCode: "FileWithSameNameExists", StatusCode: http.StatusPreconditionFailed}
	FileExists                   = &ApiError{ErrorMessage: "File with same name or id already exists under the given parent folder", ErrorCode: "FileExists", StatusCode: http.StatusPreconditionFailed}
	FileContentsLimitExceeded    = &ApiError{ErrorMessage: "Request exceeds maximum number of contents supported", StatusCode: http.StatusBadRequest, ErrorCode: "FileContentsLimitExceeded"}
	StaleObject                  = &ApiError{ErrorMessage: "Cannot update as request contains stale data", StatusCode: http.StatusBadRequest, ErrorCode: "StaleObject"}
	InvalidSkipToken             = &ApiError{ErrorMessage: "request contains invalid skip token", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidSkipToken"}
	InvalidPageSize              = &ApiError{ErrorMessage: "request contains invalid page size", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidPageSize"}
	PageSizeExceeded             = &ApiError{ErrorMessage: "The top parameter exceeds the maximum supported page size (" + strconv.FormatInt(constants.MaxPageSize, 10) + ").", StatusCode: http.StatusBadRequest, ErrorCode: "PageSizeExceeded"}
	InvalidSortParameter         = &ApiError{ErrorMessage: "given sort parameter is invalid", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidSortParameter"}
	InvalidUrl                   = &ApiError{ErrorMessage: "The url is not valid.", StatusCode: http.StatusNotFound, ErrorCode: "InvalidUrl"}
	MethodNotAllowed             = &ApiError{ErrorMessage: "Method is not allowed on this resource.", StatusCode: http.StatusMethodNotAllowed, ErrorCode: "MethodNotAllowed"}
	InvalidName                  = &ApiError{ErrorMessage: "resource name is invalid", StatusCode: http.StatusBadRequest, ErrorCode: "InvalidName"}
)
