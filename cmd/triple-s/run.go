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
		defer objectFile.Close()
		if err != nil {
			fmt.Println("3")
			return
		}

		// check if the objects.csv exists, if NOT ==> create
		_, err = os.Stat(dirPath + bucket + "/" + "objects.csv")
		if err != nil {
			if os.IsNotExist(err) {
				err := core.CreateNewObjectsCSV(dirPath, bucket, object, r, objectSize)
				if err != nil {
					fmt.Println("4")
					return
				}
			} else {
				fmt.Fprintf(os.Stderr, "Stat objects.csv stage: %s\n", err)
				return
			}
		} else {
			err := core.UpdateExistingObjMetadata(dirPath, bucket, object, r, objectSize)
			if err != nil {
				fmt.Println("5", err)
				return
			}
		}

		// UPDATE BUCKETS.CSV LastModified TIME
		err = core.UpdateExistingBucketMetadata(dirPath, bucket)
		if err != nil {
			fmt.Println("6")
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

	switch endpoint {
	case "bucket":
		if bucket == "" {

			xmlResponse, err := core.RootBucketsXML(dirPath)
			if err != nil {
				fmt.Println("10")
				return
			}

			w.Write(xmlResponse)
			// return an XML list of all bucket names and metadata
			// response with 200 OK status
		} else {
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

			xmlResponse, err := core.SingleBucketXML(dirPath, bucket)
			if err != nil {
				fmt.Println("12")
				return
			}

			w.Write(xmlResponse)
		}
	case "object":
		if object == "" {
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

			xmlResponse, err := core.BucketObjectsXML(dirPath, bucket)
			if err != nil {
				fmt.Println("11")
				return
			}
			w.Write(xmlResponse)

		} else {
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
		}
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

	switch endpoint {
	case "bucket":
		err := os.Remove(dirPath + bucket)
		if err != nil {
			fmt.Println("8")
			return
		}

		err = core.RemoveBucketMetadata(dirPath, bucket)
		if err != nil {
			fmt.Println("9")
			return
		}

	case "object":
		// CHECK IF THE BUCKET exists
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

		// remove object and metadata
		err = core.DeleteObjectAndMeta(dirPath, bucket, object)
		if err != nil {
			fmt.Println("7")
			return
		}

		err = core.UpdateExistingBucketMetadata(dirPath, bucket)
		if err != nil {
			fmt.Println("8")
			return
		}
		// coresponding response
		w.Write([]byte("Object deleted"))
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
