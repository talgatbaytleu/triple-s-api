package core

import (
	"encoding/csv"
	"os"
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

func IsBucketExist(csvfile, bucket string) bool {
	file, _ := os.Open(csvfile)
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all rows from the CSV
	records, _ := reader.ReadAll()

	// Iterate through all rows and print the first field
	for _, record := range records[1:] {
		if len(record) > 0 {
			// Access the first field (index 0)
			if record[0] == bucket {
				return true
			}
		}
	}
	return false
}

// func MetadataObjectCreation(object string, size int) string {
// 	var Size string
// 	var ContentType string
// 	var LastModified string
// }
