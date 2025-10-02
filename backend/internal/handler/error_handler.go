// Package handler provider web api handler
package handler

import (
	"blog-server/internal/dto/response"
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		code := fiber.StatusInternalServerError
		msg := "Internal Server Error"

		// 如果是 Fiber 错误
		var e *fiber.Error
		if errors.As(err, &e) {
			code = e.Code
			msg = e.Message
		}

		logger.Error(err.Error())

		if err := c.Status(code).JSON(response.Error(code, msg)); err != nil {
			return c.Status(code).SendString(msg)
		}

		return nil
	}
}
