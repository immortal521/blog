package errs

import "errors"

var (
	ErrDuplicateURL = errors.New("duplicate url")
	ErrPostNotFound = errors.New("post not found")
)
