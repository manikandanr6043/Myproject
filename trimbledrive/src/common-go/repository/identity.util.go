package repository

import (
	"strconv"
	"strings"

	"trimble.com/common/constants"
)

// GetVisibleId extracts and returns the visible id from the DB based id
func GetVisibleId(dbId *string) *string {
	if dbId != nil {
		dbIdSlice := strings.Split(*dbId, "#")
		visibleId := dbIdSlice[2]
		return &visibleId
	}
	return nil
}

// GetDbId returns DB identifier based in the visibleId, spaceId and resource type
func GetDbId(spaceId string, resourceType string, visibleId *string) *string {
	if visibleId != nil {
		dbId := spaceId + "#" + resourceType + "#" + *visibleId
		return &dbId
	}
	return nil
}

// GetDbVersionId returns DB version identifier based on id and version
func GetDbVersionId(spaceId string, resourceType string, visibleId string, majorVersion int64, minorVersion int64) *string {
	return GetDbVersionIdFromDbId(*GetDbId(spaceId, resourceType, &visibleId), majorVersion, minorVersion)
}

// GetDbVersionIdFromDbId returns DB version identifier from an existing db identifier
func GetDbVersionIdFromDbId(dbId string, majorVersion int64, minorVersion int64) *string {
	dbVersionId := dbId + "#" + strconv.FormatInt(majorVersion, 10) + constants.VersionSeparator + strconv.FormatInt(minorVersion, 10)
	return &dbVersionId
}

// GetDbVersionIdFromDbIdAndVersion returns DB version identifier from an existing db identifier and version (string)
func GetDbVersionIdFromDbIdAndVersion(dbId string, version string) *string {
	dbVersionId := dbId + "#" + version
	return &dbVersionId
}
