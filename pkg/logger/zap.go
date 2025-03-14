package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

var (
	zapLoggerInstance *ZapLogger
	mu                sync.Mutex
)

func GetZapLoggerInstance() *ZapLogger {
	if zapLoggerInstance == nil {
		mu.Lock()
		defer mu.Unlock()

		if zapLoggerInstance == nil {
			config := zap.NewDevelopmentConfig()
			config.EncoderConfig = zapcore.EncoderConfig{
				TimeKey:       "time",
				LevelKey:      "level",
				NameKey:       "logger",
				CallerKey:     "caller",
				MessageKey:    "msg",
				StacktraceKey: "stacktrace",
				LineEnding:    zapcore.DefaultLineEnding,
				EncodeTime:    zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), 
				EncodeLevel:   zapcore.CapitalLevelEncoder,
				EncodeCaller:  zapcore.ShortCallerEncoder,
			}
			logger, _ := config.Build()
			zapLoggerInstance = &ZapLogger{
				logger: logger.Sugar(),
			}
		}
	}
	return zapLoggerInstance
}

// Simple logging methods
func (l *ZapLogger) Debug(msg string, extra ...interface{}) {
	l.logger.Debugw(msg, extra...)
}

func (l *ZapLogger) Info(msg string, extra ...interface{}) {
	l.logger.Infow(msg, extra...)
}

func (l *ZapLogger) Warn(msg string, extra ...interface{}) {
	l.logger.Warnw(msg, extra...)
}

func (l *ZapLogger) Error(msg string, extra ...interface{}) {
	l.logger.Errorw(msg, extra...)
}

func (l *ZapLogger) Fatal(msg string, extra ...interface{}) {
	l.logger.Fatalw(msg, extra...)
}
