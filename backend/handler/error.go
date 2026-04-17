package handler

import (
	"blog-server/config"
	"blog-server/logger"
	"blog-server/pkg/errx"
	"blog-server/response"

	"github.com/gofiber/fiber/v3"
)

func ErrorHandler(log logger.Logger, cfg *config.Config) fiber.ErrorHandler {
	return func(c fiber.Ctx, err error) error {
		appErr := errx.ToAppError(err)
		httpCode := errx.MapToHTTPStatus(appErr.Code)

		msg := appErr.Msg
		if httpCode == 500 {
			msg = "Internal Server Error"
		}

		reqID := c.Get("X-Request-ID")

		errLogger := log.With(
			logger.String("request_id", reqID),
			logger.String("method", c.Method()),
			logger.String("remote_ip", c.IP()),
			logger.String("path", c.Path()),
			logger.String("original_url", c.OriginalURL()),
		)

		if cfg.App.Environment == config.EnvDev {
			errLogger.Error(appErr.FormatStack())
		} else {
			errLogger.Error(msg, logger.Err(err))
		}

		writeErr := c.Status(httpCode).JSON(response.Error(appErr.Code, msg))
		if writeErr != nil {
			return c.Status(httpCode).SendString(msg)
		}

		return nil
	}
}
