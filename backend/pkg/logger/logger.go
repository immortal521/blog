// Package logger
package logger

import (
	"context"
	"sync"

	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	once sync.Once // once 能确保初始化函数只被执行一次
)

type contextKey string

const loggerKey = contextKey("logger")

func Get() *zap.Logger {
	once.Do(func() {
		var err error
		log, err = zap.NewDevelopment()
		if err != nil {
			panic("初始化日志失败: " + err.Error())
		}
	})
	return log
}

// ToContext 将 logger 存入 context
func ToContext(ctx context.Context, log *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

// FromContext 从 context 中获取 logger
// 如果 context 中没有，则返回全局的 logger，保证函数总是安全的
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return l
	}
	return Get()
}
