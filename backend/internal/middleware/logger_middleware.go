// Package middleware
package middleware

import (
	"time"

	"blog-server/pkg/logger"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RequestLogger(log logger.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 为每个请求生成唯一 ID
			reqID := uuid.New().String()
			c.Response().Header().Set("X-Request-ID", reqID)

			// 创建带请求上下文的 logger
			reqLogger := log.WithFields(
				logger.Any("request_id", reqID),
				logger.String("method", c.Request().Method),
				logger.String("remote_ip", c.RealIP()),
				logger.String("path", c.Request().URL.Path),
				logger.String("user_agent", c.Request().UserAgent()),
			)

			start := time.Now()
			err := next(c) // 调用下一个处理器

			latency := time.Since(start)
			statusCode := c.Response().Status
			respSize := c.Response().Size

			fields := []logger.Field{
				logger.Int("status", statusCode),
				logger.Int64("response_size", respSize),
				logger.Duration("latency", latency),
				logger.String("original_url", c.Request().RequestURI),
				logger.String("referer", c.Request().Referer()),
			}

			if latency > 5*time.Second {
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
}
