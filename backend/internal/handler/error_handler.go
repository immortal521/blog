// Package handler provider web api handler
package handler

import (
	"blog-server/internal/dto"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	msg := "Internal Server Error"

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
		msg = e.Message
	} else {
		log.Error(err.Error())
	}

	err = c.Status(code).JSON(dto.Error(code, msg))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	return nil
}
