package handler

import (
	"blog-server/config"
	"blog-server/logger"
	"blog-server/pkg/errx"
	"blog-server/response"

	"github.com/labstack/echo/v5"
)

func ErrorHandler(
	log logger.Logger,
	cfg *config.Config,
) echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		appErr := errx.ToAppError(err)
		httpCode := errx.MapToHTTPStatus(appErr.Code)

		fields := []logger.Field{
			logger.String("method", c.Request().Method),
			logger.String("path", c.Request().URL.Path),
			logger.String("url", c.Request().URL.String()),
			logger.Int("status", httpCode),
			logger.String("ip", c.RealIP()),
		}

		if reqID := c.Request().
			Context().
			Value("request_id"); reqID != nil {
			if s, ok := reqID.(string); ok {
				fields = append(fields, logger.String("request_id", s))
			}
		}

		logCtx := log.With(fields...)

		if cfg.App.Environment == config.EnvDev {
			logCtx.Error(appErr.Error() + "\n" + appErr.StackString())
		} else {
			logCtx.Error("request failed", logger.Err(err))
		}

		publicMsg := errx.MessageForCode(appErr.Code)

		err = c.JSON(httpCode, response.Error(
			appErr.Code,
			publicMsg,
		))
		if err != nil {
			log.Error("failed to send error response", logger.Err(err))
		}
	}
}
