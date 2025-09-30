// Package handler provider web api handler
package handler

import (
	"blog-server/internal/dto/response"
	"blog-server/pkg/logger"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	log := logger.FromContext(c.UserContext())
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	// 如果是 Fiber 错误
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		msg = e.Message
	}

	log.Error(err.Error())
	if err := c.Status(code).JSON(response.Error(code, msg)); err != nil {
		return c.Status(code).SendString(msg)
	}

	return nil
}
