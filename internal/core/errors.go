package core

import "errors"

var ErrInvBucketNameLongSymbols = errors.New(
	"Bucket Name should be between 3 and 63 characters long\nOnly lowercase letters, numbers, hyphens (-), and dots (.) are allowed.",
)

var ErrInvBucketNameDashPeriod = errors.New(
	"Bcuket name must not begin or end with a hyphen and must not contain two consecutive periods or dashes.",
)

var ErrInvBucketNameIP = errors.New(
	"Bucket name must not be formatted as an IP address (e.g., 192.168.0.1).",
)

var ErrBucketAlreadyExists = errors.New("Bucket already exists, bucket name must be unique")

var (
	ErrBucketNotExist = errors.New("Bucket not exist")
	ErrObjectNotExist = errors.New("Object not exist")
)
