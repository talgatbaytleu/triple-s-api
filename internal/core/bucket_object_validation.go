package core

import (
	"regexp"
)

func ValidateBucket(bucket string) (bool, string) {
	// Bucket names must be unique across the system.

	// Names should be between 3 and 63 characters long.
	// Only lowercase letters, numbers, hyphens (-), and dots (.) are allowed.
	match, _ := regexp.MatchString("^[0-9a-z\\.\\-]{3,63}$", bucket)
	if !match {
		return false, "Bucket Name should be between 3 and 63 characters long\nOnly lowercase letters, numbers, hyphens (-), and dots (.) are allowed."
	}

	// Must not begin or end with a hyphen.
	match, _ = regexp.MatchString("^\\-|\\-$|\\-\\-|\\.\\.", bucket)
	if match {
		return false, "Bcuket name must not begin or end with a hyphen and must not contain two consecutive periods or dashes."
	}

	// Must not be formatted as an IP address (e.g., 192.168.0.1).
	match, _ = regexp.MatchString(
		"^(((1?[0-9]?[0-9])|(2[0-4][0-9])|(25[0-5]))\\.){3}((1?[0-9]?[0-9])|(2[0-4][0-9])|(25[0-5]))$",
		bucket,
	)
	if match {
		return false, "Bucket name must not be formatted as an IP address (e.g., 192.168.0.1)."
	}
	return true, ""
}
