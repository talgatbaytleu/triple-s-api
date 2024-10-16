package triples

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"triple-s/internal/core"
)

var (
	dirPath string
	port    string
)

func PutHandler(w http.ResponseWriter, r *http.Request) {
	fullPath := r.URL.Path[1:]
	endpoint := core.DetermineEndpoint(fullPath)
	bucket, object := core.SplitPath(fullPath)

	fmt.Println(endpoint)
	switch endpoint {
	case "bucket":
		isBucketNameValid, bucketNameErrMessage := core.ValidateBucket(bucket)
		if isBucketNameValid {

			err := core.CreateNewBucket(dirPath, bucket)
			if err != nil {
				fmt.Println("1")
				return
			}
			// UPDATE METADATA IN BUCKETS.CSV

			err = core.AddMetaToBucketsCSV(dirPath, bucket)
			if err != nil {
				fmt.Println("2")
				return
			}

		} else {
			http.Error(w, "400 - Bad Request "+bucketNameErrMessage, http.StatusBadRequest)
			return
		}

		// 200 RESPONSE
		w.Write([]byte("Bucket successfully created!"))

		// PUT object!!!
	case "object":
		// check if the bucket exists
		dirInfo, err := os.Stat(dirPath + bucket)
		dirInfo = dirInfo
		if err != nil {
			if os.IsNotExist(err) {
				http.Error(w, "400 - Bad Request, bucket doesn't exist", http.StatusBadRequest)
				return
			} else {
				http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
				fmt.Fprintf(os.Stderr, "Stat bucket stage: %s\n", err)
				return
			}
		}

		// create object
		objectFile, objectSize, err := core.CreateObject(dirPath, bucket, object, r)
		if err != nil {
			fmt.Println("3")
			return
		}
		// objectFile, err := os.Create(dirPath + bucket + "/" + object)
		// defer objectFile.Close()
		// if err != nil {
		// 	http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		// 	fmt.Fprintf(os.Stderr, "Object creating stage: %s\n", err)
		// 	return
		// }
		//
		// // write r.body into object
		// objectSize, err := io.Copy(objectFile, r.Body)
		// if err != nil {
		// 	http.Error(w, "500 - Internal Server Error", http.StatusInternalServerError)
		// 	fmt.Fprintf(os.Stderr, "Body content Copy to object stage: %s\n", err)
		// 	return
		// }

		// check if the objects.csv exists, if NOT ==> create
		_, err = os.Stat(dirPath + bucket + "/" + "objects.csv")
		if err != nil {
			if os.IsNotExist(err) {
				csvObjectsFile, err := os.Create(dirPath + bucket + "/" + "objects.csv")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Creating objects.csv stage: %s\n", err)
				}
				csvObjectsFile.WriteString(
					"ObjectKey,Size,ContentType,LastModified\n" + object + "," + strconv.Itoa(
						int(objectSize),
					) + "," + http.DetectContentType(
						[]byte(objectFile.Name()),
					) + "," + time.Now().
						Format("2006-01-02 15:04:05"),
				)
				defer csvObjectsFile.Close()
			} else {
				fmt.Fprintf(os.Stderr, "Stat objects.csv stage: %s\n", err)
			}
		} else {
			_, err := os.Stat(dirPath + bucket + "/" + object)

			if os.IsNotExist(err) {
				csvObjectsFile, err := os.OpenFile(dirPath+bucket+"/"+"objects.csv", os.O_APPEND|os.O_WRONLY, 0644)
				if err != nil {
					http.Error(w, "500 - Server Error", http.StatusInternalServerError)
					fmt.Fprintf(os.Stderr, "objects.csv opening stage: %s\n", err)
				}
				defer csvObjectsFile.Close()
				csvObjectsFile.WriteString("\n" + object + "," + strconv.Itoa(int(objectSize)) + "," + http.DetectContentType([]byte(objectFile.Name())) + "," + time.Now().Format("2006-01-02 15:04:05"))
			}

			// UPDATE OBJECTS.CSV IF OBJECT ALREADY EXISTS
			csvObjectsFile, err := os.Open(dirPath + bucket + "/" + "objects.csv")
			if err != nil {
				fmt.Fprintf(os.Stderr, "Updating objects.csv after creating obj if it exists: %s\n", err)
				return
			}

			csvObjectsReader := csv.NewReader((csvObjectsFile))
			defer csvObjectsFile.Close()

			csvObjectsRecords, err := csvObjectsReader.ReadAll()
			if err != nil {
				fmt.Fprintf(os.Stderr, "objects.csv parsing stage if exists: %s\n", err)
				return
			}

			for i, row := range csvObjectsRecords {
				if row[0] == object {
					csvObjectsRecords[i][1] = strconv.Itoa(int(objectSize))
					csvObjectsRecords[i][2] = http.DetectContentType([]byte(objectFile.Name()))
					csvObjectsRecords[i][3] = time.Now().Format("2006-01-02 15:04:05")
				}
			}

			csvObjectsFile, err = os.OpenFile(dirPath+bucket+"/"+"objects.csv", os.O_WRONLY, 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Opening the objects.csv to rewrite it if exists %s\n", err)
				return
			}

			csvObjectsWriter := csv.NewWriter(csvObjectsFile)

			err = csvObjectsWriter.WriteAll(csvObjectsRecords)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Updating the objects.csv if it exists : %s\n", err)
				return
			}

		}

		// UPDATE BUCKETS.CSV LastModified TIME
		csvBucketsFile, err := os.Open(dirPath + "buckets.csv")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Updating the buckets.csv after creating obj: %s\n", err)
			return
		}

		csvBucketsReader := csv.NewReader(csvBucketsFile)
		defer csvBucketsFile.Close()

		csvBucketsRecords, err := csvBucketsReader.ReadAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "buckets.csv parsing stage: %s\n", err)
			return
		}

		for i, row := range csvBucketsRecords {
			if row[0] == bucket {
				csvBucketsRecords[i][2] = time.Now().Format("2006-01-02 15:04:05")
			}
		}

		csvBucketsFile, err = os.OpenFile(dirPath+"buckets.csv", os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Opening the buckets.csv to rewrite it %s\n", err)
			return
		}

		csvBucketsWriter := csv.NewWriter(csvBucketsFile)

		err = csvBucketsWriter.WriteAll(csvBucketsRecords)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Updating the bucket last modified stage: %s\n", err)
			return
		}

		w.Write([]byte("Object successfully created!"))

	default:
		http.Error(w, "500 - Internal Server Error", 500)
		fmt.Fprintf(os.Stderr, "Endpoint in Put handler wrong")
	}
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	fullPath := r.URL.Path[1:]
	endpoint := core.DetermineEndpoint(fullPath)
	bucket, object := core.SplitPath(fullPath)

	object = object
	switch endpoint {
	case "bucket":
		if bucket == "" {
			// return an XML list of all bucket names and metadata
			// response with 200 OK status
		} else {
			// Error response
		}
	case "object":
	// check if bucket exists
	// check if file exists
	// return the binary content of the object with "content-type" set to imagee/png
	// "content-length" bytes
	default:
		fmt.Println("endpoint in Get handler wrong")
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fullPath := r.URL.Path[1:]
	endpoint := core.DetermineEndpoint(fullPath)
	bucket, object := core.SplitPath(fullPath)

	object = object
	bucket = bucket
	switch endpoint {
	case "bucket":
	// check if bucket exists and it's empty
	// remove bucket
	// update the buckets.csv
	// coresponding response
	case "object":
	// check if bucket and objects are exists
	// remove object
	// update the object.csv
	// coresponding response
	default:
		fmt.Println("endpoint in Delete handler wrong")
	}
}

func Run() {
	http.HandleFunc("PUT /", PutHandler)
	http.HandleFunc("GET /", GetHandler)
	http.HandleFunc("DELETE /", DeleteHandler)

	dirPath, port = core.InitFlags()

	// Check if buckets.csv exists. If not, it will be created
	err := core.CreateNewBucketsCSV(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Buckets.csv can NOT be created, Error: %s\n", err)
		os.Exit(1)
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
