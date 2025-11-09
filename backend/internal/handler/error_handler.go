// Package handler provider web api handler
package handler

import (
	"blog-server/internal/config"
	"blog-server/internal/response"
	"blog-server/pkg/errs"
	"blog-server/pkg/logger"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(log logger.Logger, cfg *config.Config) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		appErr := errs.ToAppError(err)
		httpCode := errs.MapToHTTPStatus(appErr.Code)
		msg := appErr.Msg

		if httpCode == 500 {
			msg = "Internal Server Error"
		}

		if cfg.App.Environment == config.EnvDev {
			log.Error(appErr.FormatStack())
		} else {
			log.Error(msg, logger.Error(err))
		}

		if err := c.Status(httpCode).JSON(response.Error(appErr.Code, msg)); err != nil {
			return c.Status(httpCode).SendString(msg)
		}

		return nil
	}
}
