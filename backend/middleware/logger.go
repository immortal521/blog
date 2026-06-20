package middleware

import (
	"time"

	"blog-server/config"
	"blog-server/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

const (
	ContextKeyRequestID = "request_id"
	ContextKeyLogger    = "logger"
	ContextKeyStart     = "request_start"
)

func RequestLogger(
	cfg *config.Config,
	log logger.Logger,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			reqID := c.Request().Header.Get(echo.HeaderXRequestID)
			if reqID == "" {
				reqID = uuid.New().String()
			}
			c.Response().Header().Set(echo.HeaderXRequestID, reqID)

			reqLogger := log.With(logger.String("request_id", reqID))
			c.Set(ContextKeyRequestID, reqID)
			c.Set(ContextKeyLogger, reqLogger)
			c.Set(ContextKeyStart, start)

			if cfg.App.IsDev() {
				reqLogger.Info("HTTP request started",
					logger.String("method", c.Request().Method),
					logger.String("path", c.Request().URL.RequestURI()),
				)
			}

			err := next(c)
			if err != nil {
				// ErrorHandler logs this request's outcome instead.
				return err
			}

			latency := time.Since(start)
			fields := []logger.Field{
				logger.String("method", c.Request().Method),
				logger.String("path", c.Request().URL.RequestURI()),
				logger.String("ip", c.RealIP()),
				logger.String("user_agent", c.Request().UserAgent()),
				logger.String("referer", c.Request().Referer()),
				logger.Duration("latency", latency),
			}

			if resp, uwErr := echo.UnwrapResponse(c.Response()); uwErr == nil {
				fields = append(fields, logger.Int("status", resp.Status))
			}

			if latency > 5*time.Second {
				fields = append(fields, logger.Bool("slow_request", true))
			}

			reqLogger.Info("HTTP request completed", fields...)
			return nil
		}
	}
}
