// Package logger
package logger

import "go.uber.org/zap"

func NewLogger() *zap.Logger {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic("logger initialization failed: " + err.Error())
	}
	return log
}
