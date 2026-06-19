package middleware

import (
	"context"
	"time"

	"blog-server/config"
	"blog-server/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

type contextKey string

const (
	requestIDKey contextKey = "request_id"
	loggerCtxKey string     = "logger"
)

// RequestLogger returns a echo handler func that logs HTTP requests and responses
// It generates a unique request ID, logs request details, and logs response with latency
func RequestLogger(
	log logger.Logger,
	cfg *config.Config,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			reqID := c.Request().Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = uuid.New().String()
			}
			c.Response().Header().Set(echo.HeaderXRequestID, reqID)

			ctx := context.WithValue(c.Request().Context(), requestIDKey, reqID)
			c.SetRequest(c.Request().WithContext(ctx))

			reqLogger := log.WithContext(ctx)
			c.Set(loggerCtxKey, reqLogger)

			if cfg.App.IsDev() {
				reqLogger.Info("HTTP request started",
					logger.String("method", c.Request().Method),
					logger.String("path", c.Request().URL.RequestURI()),
				)
			}

			err := next(c)

			latency := time.Since(start)

			fields := []logger.Field{
				logger.String("request_id", reqID),
				logger.String("method", c.Request().Method),
				logger.String("path", c.Request().URL.RequestURI()),
				logger.String("url", c.Request().URL.String()),
				logger.String("ip", c.RealIP()),
				logger.String("user_agent", c.Request().UserAgent()),
				logger.String("referer", c.Request().Referer()),
				logger.Duration("latency", latency),
			}

			if latency > 5*time.Second {
				fields = append(fields, logger.Bool("slow_request", true))
			}

			if err == nil {
				reqLogger.Info("HTTP request completed", fields...)
			}

			return err
		}
	}
}
