package logger

import (
	"blog/config"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ZapLogger struct {
	Logger *zap.SugaredLogger
}

var (
	zapLoggerInstance = &ZapLogger{}
	mu                = &sync.Mutex{}
	logLevelMapping   = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
		"fatal": zapcore.FatalLevel,
	}
)

func GetZapLoggerInstance(loggerCfg *config.Logger) Logger {
	mu.Lock()
	defer mu.Unlock()

	if zapLoggerInstance == nil {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   loggerCfg.LogFilePath,
			MaxSize:    1,
			MaxAge:     5,
			MaxBackups: 10,
			Compress:   true,
		})

		config := zap.NewProductionEncoderConfig()

		config.EncodeTime = zapcore.ISO8601TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			w,
			getLogLevel(loggerCfg.Level),
		)

		logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()

		logger = logger.With("AppName", "MyApp", "LoggerName", "ZeroLog")

		zapLoggerInstance.Logger = logger
	}
	return zapLoggerInstance
}

func getLogLevel(logLevel string) zapcore.Level {
	level, exists := logLevelMapping[logLevel]

	if !exists {
		return zapcore.DebugLevel
	}

	return level
}

func (l *ZapLogger) Debug(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)

	l.Logger.Debugw(msg, params...)
}

func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.Logger.Debugf(template, args)
}

func (l *ZapLogger) Info(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)
	l.Logger.Infow(msg, params...)
}

func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.Logger.Infof(template, args)
}

func (l *ZapLogger) Warn(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)
	l.Logger.Warnw(msg, params...)
}

func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.Logger.Warnf(template, args)
}

func (l *ZapLogger) Error(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)
	l.Logger.Errorw(msg, params...)
}

func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.Logger.Errorf(template, args)
}

func (l *ZapLogger) Fatal(cat Category, sub SubCategory, msg string, extra map[ExtraKey]interface{}) {
	params := prepareLogInfo(cat, sub, extra)
	l.Logger.Fatalw(msg, params...)
}

func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.Logger.Fatalf(template, args)
}

func prepareLogInfo(cat Category, sub SubCategory, extra map[ExtraKey]interface{}) []interface{} {
	if extra == nil {
		extra = make(map[ExtraKey]interface{})
	}

	extra["Category"] = cat

	extra["SubCategory"] = sub

	return logParamsToZapParams(extra)
}

func logParamsToZapParams(keys map[ExtraKey]interface{}) []interface{} {
	params := make([]interface{}, 0, len(keys))

	for k, v := range keys {
		params = append(params, string(k))
		params = append(params, v)
	}

	return params
}
