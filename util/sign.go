package util

import (
	"time"
)

// IsValidTimeStamp determine whether the gap between the startTime and the current timestamp is within the allowable range.
func IsValidTimeStamp(startTime int64, allowDisparity int64) bool {
	nowDisparity := time.Now().UnixMilli() - startTime

	return allowDisparity == 1 || nowDisparity <= allowDisparity
}
