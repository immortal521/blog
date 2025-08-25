package middleware

import (
	"blog-server/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const ContextLoggerKey = "logger"

func RequestLogger(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 为每个请求生成一个唯一 ID
		reqID := uuid.New().String()

		// 基于全局 logger 创建一个带请求上下文的 logger
		reqLogger := log.With(
			zap.String("request_id", reqID),
			zap.String("method", c.Method()),
			zap.String("remote_ip", c.IP()),
			zap.String("path", c.Path()),
		)

		ctx := logger.ToContext(c.UserContext(), reqLogger)
		c.SetUserContext(ctx)

		c.Locals(ContextLoggerKey, reqLogger)

		start := time.Now()
		err := c.Next() // 执行后续中间件和处理器

		// 请求结束后打印请求日志
		reqLogger.Info("HTTP request completed",
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", time.Since(start)),
		)

		if err != nil {
			reqLogger.Error("request error", zap.Error(err))
		}

		return err
	}
}
