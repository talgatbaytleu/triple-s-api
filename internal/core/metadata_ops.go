package core

import (
	"encoding/csv"
	"net/http"
	"os"
	"strconv"
	"time"
)

func MetadataBucketCreation(bucket string) string {
	var timeCreated string
	var timeModified string
	var status string

	timeCreated = time.Now().Format("2006-01-02 15:04:05")
	timeModified = timeCreated
	status = "active"

	return bucket + "," + timeCreated + "," + timeModified + "," + status + "\n"
}

func CreateNewBucketsCSV(dirPath string) error {
	_, err := os.Stat(dirPath + "buckets.csv")
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(dirPath + "buckets.csv")
			defer file.Close()
			if err != nil {
				return err
			}
			file.WriteString("Name,CreationTime,LastModifiedTime,Status\n")
		} else {
			return err
		}
	}
	return nil
}

func CreateNewObjectsCSV(
	dirPath, bucket, object string,
	r *http.Request,
	objectSize int,
) error {
	csvObjectsFile, err := os.Create(dirPath + bucket + "/" + "objects.csv")
	if err != nil {
		return err
	}

	csvObjectsFile.WriteString(
		"ObjectKey,Size,ContentType,LastModified\n" + object + "," + strconv.Itoa(
			objectSize,
		) + "," + r.Header.Get("content-type") + "," + time.Now().
			Format("2006-01-02 15:04:05") + "\n",
	)
	defer csvObjectsFile.Close()
	return nil
}

func UpdateExistingObjMetadata(
	dirPath, bucket, object string,
	r *http.Request,
	objectSize int,
) error {
	var objectAlreadyExists bool = false

	csvObjectsFile, err := os.Open(dirPath + bucket + "/" + "objects.csv")
	if err != nil {
		return err
	}

	csvObjectsReader := csv.NewReader((csvObjectsFile))
	defer csvObjectsFile.Close()

	csvObjectsRecords, err := csvObjectsReader.ReadAll()
	if err != nil {
		return err
	}

	for i, row := range csvObjectsRecords {
		if row[0] == object {
			csvObjectsRecords[i][1] = strconv.Itoa(objectSize)
			csvObjectsRecords[i][2] = r.Header.Get("content-type")
			csvObjectsRecords[i][3] = time.Now().Format("2006-01-02 15:04:05")
			objectAlreadyExists = true
		}
	}

	if !objectAlreadyExists {
		csvObjectsFile, err := os.OpenFile(
			dirPath+bucket+"/"+"objects.csv",
			os.O_APPEND|os.O_WRONLY,
			0644,
		)
		defer csvObjectsFile.Close()
		if err != nil {
			return err
		}
		csvObjectsFile.WriteString(
			object + "," + strconv.Itoa(
				int(objectSize),
			) + "," + r.Header.Get("content-type") + "," + time.Now().
				Format("2006-01-02 15:04:05") +
				"\n",
		)
	} else {
		csvObjectsFile, err = os.OpenFile(dirPath+bucket+"/"+"objects.csv", os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		csvObjectsWriter := csv.NewWriter(csvObjectsFile)

		err = csvObjectsWriter.WriteAll(csvObjectsRecords)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateExistingBucketMetadata(dirPath, bucket string) error {
	csvBucketsFile, err := os.Open(dirPath + "buckets.csv")
	if err != nil {
		return err
	}

	csvBucketsReader := csv.NewReader(csvBucketsFile)
	defer csvBucketsFile.Close()

	csvBucketsRecords, err := csvBucketsReader.ReadAll()
	if err != nil {
		return err
	}

	for i, row := range csvBucketsRecords {
		if row[0] == bucket {
			csvBucketsRecords[i][2] = time.Now().Format("2006-01-02 15:04:05")
		}
	}

	csvBucketsFile, err = os.OpenFile(dirPath+"buckets.csv", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	csvBucketsWriter := csv.NewWriter(csvBucketsFile)

	err = csvBucketsWriter.WriteAll(csvBucketsRecords)
	if err != nil {
		return err
	}

	return nil
}

func RemoveBucketMetadata(dirPath, bucket string) error {
	csvBucketsFile, err := os.Open(dirPath + "buckets.csv")
	if err != nil {
		return err
	}

	csvBucketsReader := csv.NewReader(csvBucketsFile)
	defer csvBucketsFile.Close()

	csvBucketsRecords, err := csvBucketsReader.ReadAll()
	if err != nil {
		return err
	}

	var filteredRecords [][]string
	for _, row := range csvBucketsRecords {
		if row[0] == bucket {
			continue
		}
		filteredRecords = append(filteredRecords, row)
	}

	csvBucketsFile, err = os.OpenFile(dirPath+"buckets.csv", os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	csvBucketsWriter := csv.NewWriter(csvBucketsFile)

	err = csvBucketsWriter.WriteAll(filteredRecords)
	if err != nil {
		return err
	}

	return nil
}
