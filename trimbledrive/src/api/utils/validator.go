package utils

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"trimble.com/common/api_error"
	"trimble.com/common/constants"
	"trimble.com/common/repository"
)

var rxUrlSafe = regexp.MustCompile(constants.UrlSafe)
var rxReservedName = regexp.MustCompile(constants.ReservedName)
var rxIllegalChars = regexp.MustCompile(constants.IllegalChars)

// ValidateRequestBody Validate the request body against the given object
// returns true if validation succeeds, false otherwise
func ValidateRequestBody(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			apiErr := &api_error.ApiError{
				ErrorMessage: "Please provide a valid value for: " + fieldErr.Field(),
				ErrorCode:    "InvalidPayload",
				StatusCode:   http.StatusBadRequest,
			}
			HandleApiError(c, apiErr)
			return false // exit on first error
		}
	}
	return true
}

// ValidateQueryParams Validate the request query parameters against the given object
// returns true if validation succeeds, false otherwise
func ValidateQueryParams(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		switch e := err.(type) {
		case validator.ValidationErrors:
			for _, fieldErr := range err.(validator.ValidationErrors) {
				apiErr := &api_error.ApiError{
					ErrorMessage: "Please provide a valid value for parameter: " + fieldErr.Field(),
					ErrorCode:    "InvalidQueryParam" + fieldErr.Field(),
					StatusCode:   http.StatusBadRequest,
				}
				HandleApiError(c, apiErr)
				return false // exit on first error
			}
		case *strconv.NumError:
			numError := err.(*strconv.NumError)
			apiErr := &api_error.ApiError{
				ErrorMessage: "InvalidQueryParam Value: " + numError.Num,
				ErrorCode:    "InvalidQueryParam",
				StatusCode:   http.StatusBadRequest,
			}
			HandleApiError(c, apiErr)
			return false
		default:
			apiErr := &api_error.ApiError{
				ErrorMessage: "InvalidQueryParam :" + e.Error(),
				ErrorCode:    "InvalidQueryParam",
				StatusCode:   http.StatusBadRequest,
			}
			HandleApiError(c, apiErr)
			return false
		}
	}
	return true
}

// ValidateHeaderParams Validate the request header parameters against the given object
// returns true if validation succeeds, false otherwise
func ValidateHeaderParams(c *gin.Context, obj any) bool {
	if err := c.ShouldBindHeader(obj); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			apiErr := &api_error.ApiError{
				ErrorMessage: "Please provide a valid value for parameter: " + fieldErr.Field(),
				ErrorCode:    "InvalidHeaderParam" + fieldErr.Field(),
				StatusCode:   http.StatusBadRequest,
			}
			HandleApiError(c, apiErr)
			return false // exit on first error
		}
	}
	return true
}

// ValidatePathParam Validate the request path parameter has a value or length greater than zero, or is not a space only string.
// returns value if validation succeeds, error otherwise
func ValidatePathParam(c *gin.Context, parameterName string) (string, *api_error.ApiError) {
	// Get parameter value
	paramValue := c.Param(parameterName)
	// validate if parameter value is not blank
	if len(strings.TrimSpace(paramValue)) > 0 {
		return paramValue, nil
	}
	// respond with bad request
	apiErr := &api_error.ApiError{
		ErrorMessage: "Please provide a valid value for parameter: " + parameterName,
		ErrorCode:    "InvalidParam",
		StatusCode:   http.StatusBadRequest,
	}
	return "", apiErr
}

// ValidateAndGetIfMatchHeader IfMatch header validation
func ValidateAndGetIfMatchHeader(ifMatch *string) (*repository.Version, *api_error.ApiError) {
	var ifMatchErr *api_error.ApiError
	if ifMatch != nil && *ifMatch != "*" {
		majorString, minorString, found := strings.Cut(*ifMatch, constants.VersionSeparator)
		var majorVersion int64
		var minorVersion *int64
		if found {
			if majorVersionValue, err := strconv.ParseInt(majorString, 10, 64); err == nil {
				majorVersion = majorVersionValue
			} else {
				ifMatchErr = api_error.InvalidIfMatchHeader
			}
			if minorVersionValue, err := strconv.ParseInt(minorString, 10, 64); err == nil {
				minorVersion = &minorVersionValue
			} else {
				ifMatchErr = api_error.InvalidIfMatchHeader
			}
			return &repository.Version{MajorVersion: majorVersion, MinorVersion: minorVersion}, ifMatchErr
		} else {
			ifMatchErr = api_error.InvalidIfMatchHeader
		}
	}
	return nil, ifMatchErr
}

// ValidateAndGetVersion Validate and get version
func ValidateAndGetVersion(version *string) (*repository.Version, *api_error.ApiError) {
	if version != nil {
		var versionErr *api_error.ApiError
		majorString, minorString, found := strings.Cut(*version, constants.VersionSeparator)
		var majorVersion int64
		var minorVersion *int64
		if majorVersionValue, err := strconv.ParseInt(majorString, 10, 64); err == nil {
			majorVersion = majorVersionValue
		} else {
			versionErr = api_error.InvalidVersionParam
		}
		if found && minorString != "*" {
			if minorVersionValue, err := strconv.ParseInt(minorString, 10, 64); err == nil {
				minorVersion = &minorVersionValue
			} else {
				versionErr = api_error.InvalidVersionParam
			}
		}
		return &repository.Version{MajorVersion: majorVersion, MinorVersion: minorVersion}, versionErr
	}
	return nil, nil
}

// ValidateTopParameter Validate page size
func ValidateTopParameter(top int64) *api_error.ApiError {
	if top <= 0 {
		return api_error.InvalidPageSize
	}
	if top > constants.MaxPageSize {
		return api_error.PageSizeExceeded
	}
	return nil
}

// GetFieldsQueryParam returns the comma separated fields query param as an array
func GetFieldsQueryParam(fields *[]string) *[]string {
	if fields != nil {
		var extractedFields []string
		for _, field := range *fields {
			extractedFields = append(extractedFields, strings.Split(field, ",")...)
		}
		return &extractedFields
	}
	return nil
}

// ValidateAndSetDeleteQueryParam to set the deleted query param to true for empty string
func ValidateAndSetDeleteQueryParam(c *gin.Context, deletedParam *bool) {
	if deletedParam != nil && c.Query("deleted") == "" {
		*deletedParam = true
	}
}

// ValidName Validate if the given name is valid
// Name should not contain reserved or illegal characters.
func ValidName(name *string) *api_error.ApiError {
	if name != nil {
		isValid := !rxReservedName.MatchString(*name) && !rxIllegalChars.MatchString(*name)
		if !isValid {
			return api_error.InvalidName
		}
	}
	return nil
}

// UrlSafe custom validator to validate if the given field string is url safe
var UrlSafe validator.Func = func(fl validator.FieldLevel) bool {
	fieldValueStr := fl.Field().String()
	return rxUrlSafe.MatchString(fieldValueStr)
}
