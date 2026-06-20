package handler

import (
	"time"

	"blog-server/config"
	"blog-server/logger"
	"blog-server/middleware"
	"blog-server/pkg/errx"
	"blog-server/response"

	"github.com/labstack/echo/v5"
)

func ErrorHandler(
	cfg *config.Config,
	log logger.Logger,
) echo.HTTPErrorHandler {
	return func(c *echo.Context, err error) {
		appErr := errx.ToAppError(err)
		httpCode := errx.MapToHTTPStatus(appErr.Code)

		reqLogger := log
		if l, ok := c.Get(middleware.ContextKeyLogger).(logger.Logger); ok {
			reqLogger = l
		}

		fields := []logger.Field{
			logger.String("method", c.Request().Method),
			logger.String("path", c.Request().URL.Path),
			logger.String("ip", c.RealIP()),
			logger.Int("status", httpCode),
		}

		if start, ok := c.Get(middleware.ContextKeyStart).(time.Time); ok {
			fields = append(fields, logger.Duration("latency", time.Since(start)))
		}

		logCtx := reqLogger.With(fields...)

		if cfg.App.IsDev() {
			logCtx.Error(appErr.Error() + "\n" + appErr.StackString())
		} else {
			logCtx.Error("request failed", logger.Err(err))
		}

		if sendErr := c.JSON(httpCode, response.Error(
			appErr.Code,
			errx.MessageForCode(appErr.Code),
		)); sendErr != nil {
			reqLogger.Error("failed to send error response", logger.Err(sendErr))
		}
	}
}
