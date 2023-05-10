package utils

import (
	"encoding/base64"
	"encoding/json"
	"reflect"
	"strings"

	"go.uber.org/zap"

	"trimble.com/common/api_error"
	"trimble.com/common/repository"
	"trimble.com/common/requestcontext"
)

func GenerateNextUrl(ctx *requestcontext.RequestContext, spaceId string, lastItem any, sort *repository.Sort) (*string, *api_error.ApiError) {

	// Generate skipToken (a.k.a cursor) with last returned value
	skipToken, skipTokenErr := GenerateSkipToken(ctx, spaceId, lastItem, sort)
	if skipTokenErr != nil {
		return nil, skipTokenErr
	}

	// Generate next url
	queryParams := ctx.HttpRequest().URL.Query()
	queryParams.Set("skipToken", *skipToken)
	host := GetHost(ctx.HttpRequest())
	url := "https://" + host + ctx.HttpRequest().URL.Path + "?" + queryParams.Encode()
	return &url, nil
}

// GenerateSkipToken Generate skipToken (a.k.a cursor) with last returned value
// skipToken remembers the position of the last returned item in the result set.
// skipToken is generated differently based on unique and non-unique pageable/sort fields.
// If the pageable/sort property is unique, then skipToken will be generated with lastItem's unique value
// If the pageable/sort property is not unique,
// then skipToken will be generated with combination of lastItem's unique (clock) and non-unique values.
func GenerateSkipToken(ctx *requestcontext.RequestContext, spaceId string, lastItem any, sort *repository.Sort) (*string, *api_error.ApiError) {
	skipToken := repository.Cursor{
		SpaceId: spaceId,
	}
	// Get value of sort field from last returned item
	v := reflect.ValueOf(lastItem)
	fieldValueMap := make(map[string]interface{}, len(sort.StructFields))
	for _, field := range sort.StructFields {
		f := reflect.Indirect(v).FieldByName(field)
		dbField, _ := repository.GetDBFieldName(field)
		fieldValueMap[dbField] = f.Interface()
	}

	if sort.UniquePageableProperty {
		skipToken.UniqueValue = fieldValueMap
	} else {
		// Attach clock to the cursor in case of non-unique sort field
		u := reflect.Indirect(v).FieldByName("ModifiedOnClock")
		uniqueValueMap := make(map[string]interface{}, 1)
		uniqueValueMap["modifiedOnClock"] = u.Interface()
		skipToken.UniqueValue = uniqueValueMap
		skipToken.NonUniqueValue = fieldValueMap
	}
	if marshal, err := json.Marshal(skipToken); err == nil {
		encoded := base64.URLEncoding.EncodeToString(marshal)
		return &encoded, nil
	} else {
		ctx.Logger().Error("error on marshalling skip token", zap.Error(err))
		return nil, api_error.InternalServerError
	}

}

// ValidateAndParseSkipToken Validate change and parse the token
// Return unmarshalled change token model or errors if any.
func ValidateAndParseSkipToken(ctx *requestcontext.RequestContext, spaceId string, token string) (*repository.Cursor, *api_error.ApiError) {
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		ctx.Logger().Error("Error in decoding token", zap.Error(err))
		return nil, api_error.InvalidSkipToken
	}
	var cursor *repository.Cursor
	if err = json.Unmarshal(decoded, &cursor); err != nil {
		ctx.Logger().Error("Error in unmarshalling skip token", zap.Error(err))
		return nil, api_error.InvalidSkipToken
	}
	if spaceId != cursor.SpaceId {
		ctx.Logger().Debug("Space identifier in the token doesn't match with the request", zap.Error(err))
		return nil, api_error.InvalidSkipToken
	}
	return cursor, nil
}

// ValidateAndParseSortParameter Validate, parse and get sort parameter
func ValidateAndParseSortParameter(sort *string, listVersions bool) (*repository.Sort, *api_error.ApiError) {

	// Assign default sort order
	if sort == nil {
		defaultSort := "createdOn desc"
		if listVersions {
			defaultSort = "version asc"
		}
		sort = &defaultSort
	}

	values := strings.Split(*sort, " ")
	if len(values) > 2 {
		return nil, api_error.InvalidSortParameter
	}
	if len(values) == 2 && (values[1] != "asc" && values[1] != "desc") {
		return nil, api_error.InvalidSortParameter
	}

	// Validate field name
	// If found valid, field will be mapped with db field name (Useful when db field and response field naming convention are different).
	// If no such field exists, throw err.
	structFields, dbFields := GetValidFieldName(values[0])
	if structFields == nil || dbFields == nil {
		return nil, api_error.InvalidSortParameter
	}

	// "clock" fields are considered to be unique per space
	// "name" field is unique only per parent, hence it is not recognized as unique pageable property,
	//  also treating "name" as non-unique field doesn't really harm anything.
	uniquePageableProperty := (*dbFields)[0] == "createdOnClock" || (*dbFields)[0] == "modifiedOnClock"
	if listVersions {
		// "ModifiedOnClock" or majorVersion and minorVersion are the only unique fields in "versions" object
		uniquePageableProperty = (*dbFields)[0] == "modifiedOnClock" || (*dbFields)[0] == "majorVersion"
	}

	ascending := false
	if len(values) > 1 {
		ascending = values[1] == "asc"
	}

	return &repository.Sort{DBFields: *dbFields, StructFields: *structFields, Ascending: ascending, UniquePageableProperty: uniquePageableProperty}, nil
}

// GetValidFieldName Returns mapped struct and db fields of the requested field
// Returns nil in cases of invalid field
func GetValidFieldName(field string) (*[]string, *[]string) {
	switch field {
	case "createdOn":
		mappedDBFields := []string{"createdOnClock"}
		mappedStructFields := []string{"CreatedOnClock"}
		return &mappedStructFields, &mappedDBFields
	case "modifiedOn":
		mappedDBFields := []string{"modifiedOnClock"}
		mappedStructFields := []string{"ModifiedOnClock"}
		return &mappedStructFields, &mappedDBFields
	case "name":
		mappedDBFields := []string{"name"}
		mappedStructFields := []string{"Name"}
		return &mappedStructFields, &mappedDBFields
	case "type":
		mappedDBFields := []string{"type"}
		mappedStructFields := []string{"Type"}
		return &mappedStructFields, &mappedDBFields
	case "size":
		mappedDBFields := []string{"size"}
		mappedStructFields := []string{"Size"}
		return &mappedStructFields, &mappedDBFields
	case "parentFolderId":
		mappedDBFields := []string{"parentFolderId"}
		mappedStructFields := []string{"ParentFolderId"}
		return &mappedStructFields, &mappedDBFields
	case "createdBy":
		mappedDBFields := []string{"createdBy"}
		mappedStructFields := []string{"CreatedBy"}
		return &mappedStructFields, &mappedDBFields
	case "modifiedBy":
		mappedDBFields := []string{"modifiedBy"}
		mappedStructFields := []string{"ModifiedBy"}
		return &mappedStructFields, &mappedDBFields
	case "version":
		mappedDBFields := []string{"majorVersion", "minorVersion"}
		mappedStructFields := []string{"MajorVersion", "MinorVersion"}
		return &mappedStructFields, &mappedDBFields
	case "deleted":
		mappedDBFields := []string{"deleted"}
		mappedStructFields := []string{"Deleted"}
		return &mappedStructFields, &mappedDBFields
	case "spaceId":
		mappedDBFields := []string{"spaceId"}
		mappedStructFields := []string{"SpaceId"}
		return &mappedStructFields, &mappedDBFields
	case "id":
		mappedDBFields := []string{"id"}
		mappedStructFields := []string{"Id"}
		return &mappedStructFields, &mappedDBFields
	default:
		return nil, nil
	}
}
