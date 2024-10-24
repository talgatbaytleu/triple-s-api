package core

import "errors"

var (
	ErrWrongEndpoint            = errors.New("Wrong Endpoint")
	ErrBucketNotExist           = errors.New("Bucket not exist")
	ErrObjectNotExist           = errors.New("Object not exist")
	ErrBucketNotEmpty           = errors.New("Bucket not empty")
	ErrBucketAlreadyExists      = errors.New("Bucket already exists, bucket name must be unique")
	ErrInvBucketNameLongSymbols = errors.New(
		"Bucket Name should be between 3 and 63 characters long\nOnly lowercase letters, numbers, hyphens (-), and dots (.) are allowed.",
	)
	ErrInvBucketNameDashPeriod = errors.New(
		"Bcuket name must not begin or end with a hyphen and must not contain two consecutive periods or dashes.",
	)
	ErrInvBucketNameIP = errors.New(
		"Bucket name must not be formatted as an IP address (e.g., 192.168.0.1).",
	)
)
