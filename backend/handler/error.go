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

		reqID := c.Get("X-Request-ID")

		logCtx := log.With(
			logger.String("request_id", reqID),
			logger.String("method", c.Method()),
			logger.String("path", c.Path()),
		)

		if cfg.App.Environment == config.EnvDev {
			logCtx.Error(appErr.Error() + "\n" + appErr.StackString())
		} else {
			logCtx.Error("request failed", logger.Err(err))
		}

		publicMsg := errx.MessageForCode(appErr.Code)

		return c.Status(httpCode).JSON(response.Error(
			appErr.Code,
			publicMsg,
		))
	}
}
