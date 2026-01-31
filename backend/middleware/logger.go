// Package middleware provides HTTP middleware for the blog system
package middleware

import (
	"time"

	"blog-server/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestLogger returns a Fiber handler that logs HTTP requests and responses
// It generates a unique request ID, logs request details, and logs response with latency
func RequestLogger(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := uuid.New().String()

		c.Set("X-Request-ID", reqID)

		reqLogger := log.WithFields(
			logger.Any("request_id", reqID),
			logger.String("method", c.Method()),
			logger.String("remote_ip", c.IP()),
			logger.String("path", c.Path()),
			logger.String("user_agent", string(c.Request().Header.UserAgent())),
		)

		start := time.Now()
		err := c.Next()

		latency := time.Since(start)

		statusCode := c.Response().StatusCode()
		respSize := c.Response().Header.ContentLength()

		fields := []logger.Field{
			logger.Int("status", statusCode),
			logger.Int("response_size", respSize),
			logger.Duration("latency", latency),
			logger.String("original_url", c.OriginalURL()),
			logger.String("referer", string(c.Request().Header.Referer())),
		}

		// Flag slow requests (>5 seconds)
		if latency > time.Second*5 {
			fields = append(fields, logger.Bool("slow_request", true))
		}

		switch {
		case statusCode >= 500:
			reqLogger.Error("HTTP request failed", fields...)
		case statusCode >= 400:
			reqLogger.Warn("HTTP client error", fields...)
		default:
			reqLogger.Info("HTTP request completed", fields...)
		}

		if err != nil {
			reqLogger.Error("Request processing error", logger.Error(err))
		}

		return err
	}
}
