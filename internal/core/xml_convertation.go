package core

import (
	"encoding/csv"
	"encoding/xml"
	"net/http"
	"os"
)

type Object struct {
	ObjectKey        string
	ObjectsSize      string
	ContentType      string
	LastModifiedTime string
}

type bucketContent struct {
	XMLName xml.Name `xml:"BucketContent"`
	Objects []Object `xml:"Object"`
}

type Bucket struct {
	Name             string
	CreationTime     string `xml:"CreationTime"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string
}

type rootContent struct {
	XMLName xml.Name `xml:"Buckets"`
	Buckets []Bucket `xml:"Bucket"`
}

type errStruct struct {
	XMLName      xml.Name `xml:"Error"`
	StatusCode   int      `xml:"Code"`
	ErrorMessage string   `xml:"Message"`
	Resource     string
}

func RootBucketsXML(dirPath string) ([]byte, error) {
	file, err := os.Open(dirPath + "buckets.csv")
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	csvReader := csv.NewReader(file)

	csvRecords, err := csvReader.ReadAll()
	if err != nil {
		return make([]byte, 0), err
	}

	var buckets []Bucket

	for i, row := range csvRecords {
		if i == 0 {
			continue
		}

		bucket := Bucket{
			Name:             row[0],
			CreationTime:     row[1],
			LastModifiedTime: row[2],
			Status:           row[3],
		}

		buckets = append(buckets, bucket)
	}

	rootBuckets := rootContent{Buckets: buckets}

	xmlData, err := xml.MarshalIndent(rootBuckets, "", " ")
	if err != nil {
		return make([]byte, 0), nil
	}

	return xmlData, nil
}

func SingleBucketXML(dirPath, bucketName string) ([]byte, error) {
	file, err := os.Open(dirPath + "buckets.csv")
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	csvReader := csv.NewReader(file)

	csvRecords, err := csvReader.ReadAll()
	if err != nil {
		return make([]byte, 0), err
	}

	var singleBucket Bucket

	for i, row := range csvRecords {
		if i == 0 || row[0] != bucketName {
			continue
		}

		singleBucket.Name = row[0]
		singleBucket.CreationTime = row[1]
		singleBucket.LastModifiedTime = row[2]
		singleBucket.Status = row[3]

	}

	xmlData, err := xml.MarshalIndent(singleBucket, "", " ")
	if err != nil {
		return make([]byte, 0), nil
	}

	return xmlData, nil
}

func BucketObjectsXML(dirPath, bucketName string) ([]byte, error) {
	file, err := os.Open(dirPath + bucketName + "/" + "objects.csv")
	defer file.Close()
	if err != nil {
		return make([]byte, 0), err
	}

	csvReader := csv.NewReader(file)

	csvRecords, err := csvReader.ReadAll()
	if err != nil {
		return make([]byte, 0), err
	}

	var objects []Object

	for i, row := range csvRecords {
		if i == 0 {
			continue
		}

		object := Object{
			ObjectKey:        row[0],
			ObjectsSize:      row[1],
			ContentType:      row[2],
			LastModifiedTime: row[3],
		}

		objects = append(objects, object)
	}

	bucketObjects := bucketContent{Objects: objects}

	xmlData, err := xml.MarshalIndent(bucketObjects, "", " ")
	if err != nil {
		return make([]byte, 0), nil
	}

	return xmlData, nil
}

func ResponseErrorXML(err error, urlpath string, w http.ResponseWriter) {
	var errstr errStruct
	errstr.Resource = urlpath
	errstr.ErrorMessage = err.Error()

	switch err {
	case ErrBucketNotExist, ErrObjectNotExist:
		errstr.StatusCode = 404
	case ErrBucketAlreadyExists, ErrBucketNotEmpty:
		errstr.StatusCode = 409
	case ErrInvBucketNameIP,
		ErrInvBucketNameDashPeriod,
		ErrInvBucketNameLongSymbols,
		ErrWrongEndpoint:
		errstr.StatusCode = 400
	default:
		errstr.StatusCode = 500
	}

	xmlData, _ := xml.MarshalIndent(errstr, "", " ")
	http.Error(w, string(xmlData), errstr.StatusCode)
	return
}
