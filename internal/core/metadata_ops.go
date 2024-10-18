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

	return "\n" + bucket + "," + timeCreated + "," + timeModified + "," + status
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
			file.WriteString("Name,CreationTime,LastModifiedTime,Status")
		} else {
			return err
		}
	}
	return nil
}

func CreateNewObjectsCSV(
	dirPath, bucket, object string,
	objectSize int,
	objectFile *os.File,
) error {
	csvObjectsFile, err := os.Create(dirPath + bucket + "/" + "objects.csv")
	if err != nil {
		return err
	}
	csvObjectsFile.WriteString(
		"ObjectKey,Size,ContentType,LastModified\n" + object + "," + strconv.Itoa(
			objectSize,
		) + "," + http.DetectContentType(
			[]byte(objectFile.Name()),
		) + "," + time.Now().
			Format("2006-01-02 15:04:05") + "\n",
	)
	defer csvObjectsFile.Close()
	return nil
}

func UpdateExistingObjMetadata(
	dirPath, bucket, object string,
	objectSize int,
	objectFile *os.File,
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
			csvObjectsRecords[i][1] = strconv.Itoa(int(objectSize))
			csvObjectsRecords[i][2] = http.DetectContentType([]byte(objectFile.Name()))
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
			) + "," + http.DetectContentType(
				[]byte(objectFile.Name()),
			) + "," + time.Now().
				Format("2006-01-02 15:04:05") +
				"\n",
		)
	} else {
		csvObjectsFile, err = os.OpenFile(dirPath+bucket+"/"+"objects.csv", os.O_WRONLY, 0644)
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
