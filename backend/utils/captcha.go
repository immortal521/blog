// Package utils
package utils

func GenerateCaptcha() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return RandomString(6, charset)
}
