package core

import (
	"io"
	"net/http"
	"os"
)

func CreateObject(dirPath, bucket, object string, r *http.Request) (*os.File, int, error) {
	objectFile, err := os.Create(dirPath + bucket + "/" + object)
	defer objectFile.Close()
	if err != nil {
		return nil, 0, err
	}
	objectSize, err := io.Copy(objectFile, r.Body)
	if err != nil {
		return nil, 0, err
	}

	return objectFile, int(objectSize), nil
}
