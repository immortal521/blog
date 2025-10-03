package logger

import (
	"blog-server/internal/config"
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	l *zap.Logger
}

func NewLogger(cfg *config.Config) Logger {
	var encoder zapcore.Encoder
	var level zapcore.Level

	if cfg.App.Environment == "production" {
		// 生产环境：JSON 结构化日志
		encoderCfg := zap.NewProductionEncoderConfig()
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
		encoder = zapcore.NewJSONEncoder(encoderCfg)
		level = zap.InfoLevel
	} else {
		// 开发环境：彩色控制台日志
		encoderCfg := zapcore.EncoderConfig{
			TimeKey:       "time",
			LevelKey:      "level",
			CallerKey:     "caller",
			MessageKey:    "msg",
			StacktraceKey: "stack_trace",

			EncodeDuration: zapcore.StringDurationEncoder,

			EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString("[" + t.Format("2006-01-02 15:04:05 -0700") + "]")
			},

			EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
				var color string
				switch l {
				case zapcore.DebugLevel, zapcore.InfoLevel:
					color = "\033[32m" // 绿色
				case zapcore.WarnLevel:
					color = "\033[33m" // 黄色
				case zapcore.ErrorLevel, zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
					color = "\033[31m" // 红色
				default:
					color = "\033[0m"
				}
				enc.AppendString(fmt.Sprintf("%s[%s]%s", color, l.CapitalString(), "\033[0m"))
			},

			EncodeCaller: func(ec zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString("[" + ec.TrimmedPath() + "]")
			},
		}
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		level = zap.DebugLevel
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(zapcore.Lock(os.Stdout)),
		level,
	)

	l := zap.New(core, zap.AddCaller())
	return &zapLogger{l: l}
}

func (z *zapLogger) Fatal(msg string, fields ...Field) {
	z.l.Fatal(msg, convertFields(fields)...)
}

func (z *zapLogger) Panic(msg string, fields ...Field) {
	z.l.Panic(msg, convertFields(fields)...)
}

func (z *zapLogger) WithFields(fields ...Field) Logger {
	return &zapLogger{l: z.l.With(convertFields(fields)...)}
}

func (z *zapLogger) Debug(msg string, fields ...Field) {
	z.l.Info(msg, convertFields(fields)...)
}

func (z *zapLogger) Error(msg string, fields ...Field) {
	z.l.Error(msg, convertFields(fields)...)
}

func (z *zapLogger) Info(msg string, fields ...Field) {
	z.l.Info(msg, convertFields(fields)...)
}

func (z *zapLogger) Sync() error {
	return z.l.Sync()
}

func (z *zapLogger) Warn(msg string, fields ...Field) {
	z.l.Warn(msg, convertFields(fields)...)
}

func convertFields(fields []Field) []zap.Field {
	zf := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		switch f.Type {
		case fieldError:
			zf = append(zf, zap.Error(f.Err))
		default:
			zf = append(zf, zap.Any(f.Key, f.Any))
		}
	}
	return zf
}
