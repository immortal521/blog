// Package errs
package errs

import "github.com/gofiber/fiber/v2"

// BadRequest 400 Bad Request
func BadRequest(msg string) error {
	return fiber.NewError(fiber.StatusBadRequest, msg)
}

// Unauthorized 401 Unauthorized
func Unauthorized(msg string) error {
	return fiber.NewError(fiber.StatusUnauthorized, msg)
}

// Forbidden 403 Forbidden
func Forbidden(msg string) error {
	return fiber.NewError(fiber.StatusForbidden, msg)
}

// NotFound 404 Not Found
func NotFound(msg string) error {
	return fiber.NewError(fiber.StatusNotFound, msg)
}

// Conflict 409 Conflict
func Conflict(msg string) error {
	return fiber.NewError(fiber.StatusConflict, msg)
}

// UnprocessableEntity 422 Unprocessable Entity
func UnprocessableEntity(msg string) error {
	return fiber.NewError(fiber.StatusUnprocessableEntity, msg)
}

// InternalServerError 500 Internal Server Error
func InternalServerError(msg string) error {
	return fiber.NewError(fiber.StatusInternalServerError, msg)
}

// BadGateway 502 Bad Gateway
func BadGateway(msg string) error {
	return fiber.NewError(fiber.StatusBadGateway, msg)
}

// ServiceUnavailable 503 Service Unavailable
func ServiceUnavailable(msg string) error {
	return fiber.NewError(fiber.StatusServiceUnavailable, msg)
}
