package core

import (
	"strings"
)

func SplitPath(fullPath string) (string, string) {
	var bucket, object string

	SlicePath := strings.Split(fullPath, "/")
	bucket = SlicePath[0]

	if len(SlicePath) > 1 {
		object = SlicePath[1]
	}

	return bucket, object
}

func DetermineEndpoint(fullPath string) string {
	lenOfSlice := len(strings.Split(fullPath, "/"))

	if lenOfSlice == 2 {
		return "object"
	} else if lenOfSlice == 1 {
		return "bucket"
	}
	return "unknown"
}
