package core

import (
	"strings"
)

func ValidatePath(fullPath string) (string, string, error) {
	var bucket, object string
	var err error

	BucObjSlice := strings.Split(fullPath, "/")

	if len(BucObjSlice) > 0 {
		if len(BucObjSlice) > 1 {
			object = BucObjSlice[1]
		}
		bucket = BucObjSlice[0]
	} else {
		// err = Err
		// ERROR MUST BE HANDLED HERE!!!
		return "", "", err
	}

	return bucket, object, nil
}

func DetermineEndpoint(fullPath string) string {
	lenOfSlice := len(strings.Split(fullPath, "/ "))

	if lenOfSlice > 1 {
		return "object"
	}
	return "bucket"
}

func ValidateBucket(bucket string) bool {
	// Bucket names must be unique across the system.

	// Names should be between 3 and 63 characters long.
	if len(bucket) < 3 || len(bucket) > 63 {
		return false
	}
	// Only lowercase letters, numbers, hyphens (-), and dots (.) are allowed.
	if !IsLowerAlfaNumeric(bucket) {
		return false
	}
	// Must not begin or end with a hyphen and must not contain two consecutive periods or dashes.
	if bucket[0] == '-' || strings.Contains(bucket, "--") || strings.Contains(bucket, "..") {
		return false
	}
	return true
}

// is bucket name lowercase alfanumeric and contains only -(dashes) and . (periods)
func IsLowerAlfaNumeric(bucket string) bool {
	for i := range bucket {
		if i < 45 || i == 47 || (i > 57 && i < 97) || i > 122 {
			return false
		}
	}
	return true
}

func IsNotIP(bucket string) bool {
	return true
}
