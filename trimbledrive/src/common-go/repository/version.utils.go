package repository

import (
	"strconv"

	"trimble.com/common/constants"
)

// FormatVersion returns version as float64 from the given major and minor version
func FormatVersion(majorVersion int64, minorVersion int64) string {
	return strconv.FormatInt(majorVersion, 10) + constants.VersionSeparator + strconv.FormatInt(minorVersion, 10)
}
