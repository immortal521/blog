// Package handler provider web api handler
package handler

import (
	"blog-server/internal/dto/response"
	"blog-server/pkg/errs"
	"blog-server/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(log logger.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		appErr := errs.ToAppError(err)
		httpCode := errs.MapToHTTPStatus(appErr.Code)

		log.Error(appErr.FormatStack())

		if err := c.Status(httpCode).JSON(response.Error(appErr.Code, appErr.Msg)); err != nil {
			return c.Status(httpCode).SendString(appErr.Msg)
		}

		return nil
	}
}
