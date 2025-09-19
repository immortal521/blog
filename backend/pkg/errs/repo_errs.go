package errs

import "errors"

var (
	ErrDuplicateURL = errors.New("duplicate url")
	ErrPostNotFound = errors.New("post not found")

	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user exists")
)
