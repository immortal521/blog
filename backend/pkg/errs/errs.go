// Package errs
package errs

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/gofiber/fiber/v2"
)

type AppError struct {
	Code  int
	Msg   string
	Err   error
	stack []uintptr
}

func New(code int, msg string, err error) *AppError {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(2, pcs)
	pcs = pcs[:n]
	return &AppError{
		Code:  code,
		Msg:   msg,
		Err:   err,
		stack: pcs,
	}
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %+v", e.Code, e.Msg, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func (e *AppError) StackString() string {
	if len(e.stack) == 0 {
		return ""
	}
	frames := runtime.CallersFrames(e.stack)
	var s string
	for {
		frames, more := frames.Next()
		s += fmt.Sprintf("%s\n\t%s:%d\n", frames.Function, frames.File, frames.Line)
		if !more {
			break
		}
	}
	return s
}

func (e *AppError) FormatStack() string {
	if e == nil {
		return ""
	}
	s := e.Error()
	stackStr := e.StackString()
	if stackStr != "" {
		s += "\n Stack trace:\n" + stackStr
	}
	if e.Err != nil {
		if nested, ok := e.Err.(*AppError); ok {
			s += "\nCaused by:\n" + nested.FormatStack()
		}
	}
	return s
}

func ToAppError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return New(CodeInternalError, "Internal Server Error", err)
}

// NoContent 204 No Content
func NoContent(msg string) error {
	return fiber.NewError(fiber.StatusNoContent, msg)
}

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
