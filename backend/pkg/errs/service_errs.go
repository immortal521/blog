package errs

import "errors"

var (
	ErrInvalidCaptcha = errors.New("invalid captcha")
	ErrInvalidToken   = errors.New("invalid token")
)
