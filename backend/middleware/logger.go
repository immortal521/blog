// Package middleware provides HTTP middleware for the blog system
package middleware

import (
	"time"

	"blog-server/config"
	"blog-server/errs"
	"blog-server/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequestLogger returns a Fiber handler that logs HTTP requests and responses
// It generates a unique request ID, logs request details, and logs response with latency
func RequestLogger(log logger.Logger, cfg *config.Config) fiber.Handler {
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

		if cfg.App.Environment == config.EnvDev {
			reqLogger.Info("HTTP request started")
		}

		start := time.Now()
		err := c.Next()
		latency := time.Since(start)

		statusCode := c.Response().StatusCode()
		// 用错误映射推导最终 status，保证和 ErrorHandler 一致
		if err != nil {
			appErr := errs.ToAppError(err)
			statusCode = errs.MapToHTTPStatus(appErr.Code)
		}

		respSize := c.Response().Header.ContentLength()

		fields := []logger.Field{
			logger.Int("status", statusCode),
			logger.Int("response_size", respSize),
			logger.Duration("latency", latency),
			logger.String("original_url", c.OriginalURL()),
			logger.String("referer", string(c.Request().Header.Referer())),
		}

		if latency > 5*time.Second {
			fields = append(fields, logger.Bool("slow_request", true))
		}

		// 只记录“结果日志”，不再单独打 err 详情
		switch {
		case statusCode >= 500:
			reqLogger.Error("HTTP request completed", fields...)
		case statusCode >= 400:
			reqLogger.Warn("HTTP request completed", fields...)
		default:
			reqLogger.Info("HTTP request completed", fields...)
		}

		return err
	}
}
