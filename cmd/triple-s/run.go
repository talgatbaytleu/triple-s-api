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
		err := core.ValidateBucket(bucket)
		if err == nil {

			err := core.CreateNewBucket(dirPath, bucket)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}
			// UPDATE METADATA IN BUCKETS.CSV

			err = core.AddMetaToBucketsCSV(dirPath, bucket)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}

		} else {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}

		// 200 RESPONSE
		w.Write([]byte("Bucket successfully created!"))
		return

		// PUT object!!!
	case "object":
		// check if the bucket exists

		err := core.CheckBucketExist(dirPath, bucket)
		if err != nil {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}
		// create object
		objectFile, objectSize, err := core.CreateObject(dirPath, bucket, object, r)
		defer objectFile.Close()
		if err != nil {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}

		// check if the objects.csv exists, if NOT ==> create
		_, err = os.Stat(dirPath + bucket + "/" + "objects.csv")
		if err != nil {
			if os.IsNotExist(err) {
				err := core.CreateNewObjectsCSV(dirPath, bucket, object, r, objectSize)
				if err != nil {
					core.ResponseErrorXML(err, fullPath, w)
					return
				}
			} else {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}
		} else {
			err := core.UpdateExistingObjMetadata(dirPath, bucket, object, r, objectSize)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}
		}

		// UPDATE BUCKETS.CSV LastModified TIME
		err = core.UpdateExistingBucketMetadata(dirPath, bucket)
		if err != nil {
			fmt.Println("6", err)
			return
		}
		w.Write([]byte("Object successfully created!"))
		return

	default:
		core.ResponseErrorXML(core.ErrWrongEndpoint, fullPath, w)
		return
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
				core.ResponseErrorXML(err, fullPath, w)
				return
			}

			w.Write(xmlResponse)
			return
		} else {
			err := core.CheckBucketExist(dirPath, bucket)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}

			xmlResponse, err := core.SingleBucketXML(dirPath, bucket)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}

			w.Write(xmlResponse)
			return
		}
	case "object":
		if object == "" {

			err := core.CheckBucketExist(dirPath, bucket)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}
			xmlResponse, err := core.BucketObjectsXML(dirPath, bucket)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}
			w.Write(xmlResponse)

		} else {
			err := core.CheckBucketExist(dirPath, bucket)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}

			err = core.CheckObjectExist(dirPath, bucket, object)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}

			objContent, err := os.ReadFile(dirPath + bucket + "/" + object)
			if err != nil {
				core.ResponseErrorXML(err, fullPath, w)
				return
			}
			w.Write(objContent)
			return
		}
	default:
		core.ResponseErrorXML(core.ErrWrongEndpoint, fullPath, w)
		return
	}
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	fullPath := r.URL.Path[1:]
	endpoint := core.DetermineEndpoint(fullPath)
	bucket, object := core.SplitPath(fullPath)

	switch endpoint {
	case "bucket":
		err := core.CheckBukcetEmpty(dirPath, bucket)
		if err != nil {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}

		err = os.Remove(dirPath + bucket)
		if err != nil {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}

		err = core.RemoveBucketMetadata(dirPath, bucket)
		if err != nil {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}

		w.WriteHeader(204)
		return
	case "object":
		// CHECK IF THE BUCKET exists

		err := core.CheckBucketExist(dirPath, bucket)
		if err != nil {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}

		// remove object and metadata (object existens checked inside)
		err = core.DeleteObjectAndMeta(dirPath, bucket, object)
		if err != nil {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}

		err = core.UpdateExistingBucketMetadata(dirPath, bucket)
		if err != nil {
			core.ResponseErrorXML(err, fullPath, w)
			return
		}
		// coresponding response
		w.WriteHeader(204)
		return
	default:
		core.ResponseErrorXML(core.ErrWrongEndpoint, fullPath, w)
		return
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
