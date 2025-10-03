// Package middleware
package middleware

import (
	"blog-server/pkg/logger"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const ContextLoggerKey = "logger"

func RequestLogger(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 为每个请求生成一个唯一 ID
		reqID := uuid.New().String()

		fmt.Println("RemoteAddr:", c.Context().RemoteAddr())
		fmt.Println("X-Forwarded-For:", c.Get("X-Forwarded-For"))
		fmt.Println("X-Real-IP:", c.Get("X-Real-IP"))
		fmt.Println("c.IP():", c.IP())

		// 基于全局 logger 创建一个带请求上下文的 logger
		reqLogger := log.WithFields(
			logger.Any("request_id", reqID),
			logger.String("method", c.Method()),
			logger.String("remote_ip", c.IP()),
			logger.String("path", c.Path()),
			logger.String("user_agent", string(c.Request().Header.UserAgent())),
		)

		c.Locals(ContextLoggerKey, reqLogger)

		start := time.Now()
		err := c.Next() // 执行后续中间件和处理器

		// 计算请求处理时间
		latency := time.Since(start)

		// 根据状态码确定日志级别
		statusCode := c.Response().StatusCode()
		fields := []logger.Field{
			logger.Int("status", statusCode),
			logger.Duration("latency", latency),
			logger.String("referer", string(c.Request().Header.Referer())),
		}

		// 对于慢请求添加警告
		if latency > time.Second*5 {
			fields = append(fields, logger.Bool("slow_request", true))
		}

		// 根据状态码记录不同级别的日志
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

// RateLimiter 速率限制中间件
func RateLimiter(rc cache.RedisClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 使用IP地址作为限流键
		key := fmt.Sprintf("rate_limit:%s", c.IP())

		// 限制每个IP每分钟最多100个请求
		limit := 100
		window := time.Minute

		// 使用Redis实现滑动窗口限流
		ctx := c.UserContext()
		now := time.Now().Unix()

		// 清除窗口外的请求记录
		err := rc.Raw().ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-int64(window.Seconds()))).Err()
		if err != nil {
			// 如果Redis出错，记录日志但不阻断请求
			if logger, ok := c.Locals(ContextLoggerKey).(*zap.Logger); ok {
				logger.Warn("速率限制Redis清理失败", zap.Error(err))
			}
		}

		// 获取当前窗口内的请求数
		count, err := rc.Raw().ZCard(ctx, key).Result()
		if err != nil {
			// 如果Redis出错，记录日志但不阻断请求
			if logger, ok := c.Locals(ContextLoggerKey).(*zap.Logger); ok {
				logger.Warn("速率限制Redis计数失败", zap.Error(err))
			}
		} else if count >= int64(limit) {
			// 超过限制，返回429状态码
			return fiber.ErrTooManyRequests
		}

		// 记录当前请求
		err = rc.Raw().ZAdd(ctx, key, redis.Z{
			Score:  float64(now),
			Member: now,
		}).Err()
		if err != nil {
			// 如果Redis出错，记录日志但不阻断请求
			if logger, ok := c.Locals(ContextLoggerKey).(*zap.Logger); ok {
				logger.Warn("速率限制Redis记录失败", zap.Error(err))
			}
		}

		// 设置过期时间以自动清理
		err = rc.Raw().Expire(ctx, key, window).Err()
		if err != nil {
			// 如果Redis出错，记录日志但不阻断请求
			if logger, ok := c.Locals(ContextLoggerKey).(*zap.Logger); ok {
				logger.Warn("速率限制Redis设置过期时间失败", zap.Error(err))
			}
		}

		return c.Next()
	}
}
