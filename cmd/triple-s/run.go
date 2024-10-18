package triples

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
		_, err := os.Stat(dirPath + bucket)
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

		// check if the objects.csv exists, if NOT ==> create
		_, err = os.Stat(dirPath + bucket + "/" + "objects.csv")
		if err != nil {
			if os.IsNotExist(err) {
				err := core.CreateNewObjectsCSV(dirPath, bucket, object, objectSize, objectFile)
				if err != nil {
					fmt.Println("4")
					return
				}
			} else {
				fmt.Fprintf(os.Stderr, "Stat objects.csv stage: %s\n", err)
				return
			}
		} else {
			err := core.UpdateExistingObjMetadata(dirPath, bucket, object, objectSize, objectFile)
			if err != nil {
				fmt.Println("5")
				return
			}
		}

		// UPDATE BUCKETS.CSV LastModified TIME
		err = core.UpdateExistingBucketMetadata(dirPath, bucket)
		if err != nil {
			fmt.Println("6")
			return
		}
		// csvBucketsFile, err := os.Open(dirPath + "buckets.csv")
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Updating the buckets.csv after creating obj: %s\n", err)
		// 	return
		// }
		//
		// csvBucketsReader := csv.NewReader(csvBucketsFile)
		// defer csvBucketsFile.Close()
		//
		// csvBucketsRecords, err := csvBucketsReader.ReadAll()
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "buckets.csv parsing stage: %s\n", err)
		// 	return
		// }
		//
		// for i, row := range csvBucketsRecords {
		// 	if row[0] == bucket {
		// 		csvBucketsRecords[i][2] = time.Now().Format("2006-01-02 15:04:05")
		// 	}
		// }
		//
		// csvBucketsFile, err = os.OpenFile(dirPath+"buckets.csv", os.O_WRONLY, 0644)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Opening the buckets.csv to rewrite it %s\n", err)
		// 	return
		// }
		//
		// csvBucketsWriter := csv.NewWriter(csvBucketsFile)
		//
		// err = csvBucketsWriter.WriteAll(csvBucketsRecords)
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Updating the bucket last modified stage: %s\n", err)
		// 	return
		// }

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
