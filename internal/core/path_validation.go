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
