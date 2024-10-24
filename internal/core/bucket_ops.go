package core

import (
	"io"
	"os"
)

func CreateNewBucket(dirPath, bucket string) error {
	err := os.Mkdir(dirPath+bucket, 0750)
	if err != nil {
		if os.IsExist(err) {
			return err
		} else {
			return err
		}
	}

	return nil
}

func AddMetaToBucketsCSV(dirPath, bucket string) error {
	csvBucketsFile, err := os.OpenFile(dirPath+"buckets.csv", os.O_APPEND|os.O_WRONLY, 0644)
	defer csvBucketsFile.Close()
	if err != nil {
		return err
	}

	csvBucketsFile.WriteString(MetadataBucketCreation(bucket))
	return nil
}

func CheckBucketExist(dirPath, bucket string) error {
	_, err := os.Stat(dirPath + bucket)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrBucketNotExist
		} else {
			return err
		}
	}
	return nil
}

func CheckBukcetEmpty(dirPath, bucket string) error {
	dir, err := os.Open(dirPath + bucket)
	defer dir.Close()
	if err != nil {
		return err
	}

	_, err = dir.Readdir(1)
	if err == io.EOF {
		return ErrBucketNotEmpty
	}

	return nil
}
