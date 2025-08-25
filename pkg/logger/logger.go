package logger

import (
	"go.uber.org/zap"
	"sync"
)

var (
	log  *zap.Logger
	once sync.Once // once 能确保初始化函数只被执行一次
)

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
