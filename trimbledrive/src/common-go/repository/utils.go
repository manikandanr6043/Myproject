package repository

import (
	"reflect"
	"strings"
)

type Cursor struct {
	SpaceId        string
	UniqueValue    map[string]interface{}
	NonUniqueValue map[string]interface{}
}

type Sort struct {
	DBFields               []string
	StructFields           []string
	Ascending              bool
	UniquePageableProperty bool
}

// GetDBSortDirectionAndRangeOperator Get DB query equivalent of sort direction and range operators
// returns 1, $gt for asc, -1, $lt for desc
func GetDBSortDirectionAndRangeOperator(sort Sort) (int, string) {
	if sort.Ascending {
		return 1, "$gt"
	} else {
		return -1, "$lt"
	}
}

// GetDBFieldName Returns db field name of the given struct field
// This can be consumed by "latest" or "versions" collection queries.
func GetDBFieldName(fieldName string) (string, bool) {
	t := reflect.TypeOf(Latest{})
	sf, ok := t.FieldByName(fieldName)
	if !ok {
		return "", false
	}
	name, tagOk := sf.Tag.Lookup("bson")
	if !tagOk {
		return "", false
	}
	return strings.Split(name, ",")[0], true
}
