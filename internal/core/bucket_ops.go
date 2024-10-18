package core

import (
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
