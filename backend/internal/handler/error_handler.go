// Package handler provider web api handler
package handler

import (
	"blog-server/internal/dto"
	"blog-server/pkg/errs"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	switch {
	case errors.Is(err, errs.ErrPostNotFound):
		code = fiber.StatusNotFound
		msg = "post not found"
	case errors.Is(err, errs.ErrDuplicateURL):
		code = fiber.StatusConflict
		msg = "url is duplicated"
	default:
		// 如果是 Fiber 错误
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			msg = e.Message
		} else {
			// 未知错误
			code = fiber.StatusInternalServerError
			msg = "Internal Server Error"
			log.Error(err.Error()) // 可以扩展打印堆栈
		}
	}

	if err := c.Status(code).JSON(dto.Error(code, msg)); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
