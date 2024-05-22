package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}
var MyLogger = ZapLogger{}
func InitZapLogger()  {
	if _,err := os.Stat("pkg/logger/log.log");err == nil{
		err = os.Truncate("pkg/logger/log.log",0)
		if err != nil{
			panic(err)
		}
	}
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"pkg/logger/log.log"}
	cfg.Encoding = "json"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	l, err := cfg.Build()
	if err != nil {
		panic("cannot initi  alize logger: " + err.Error())
	}
	MyLogger.logger = l
}

func (zl *ZapLogger) Info(message string, fields map[string]interface{}) {
	zl.logger.Info(message, convertToZapFields(fields)...)
}

func (zl *ZapLogger) Error(message string, fields map[string]interface{}) {
	zl.logger.Error(message, convertToZapFields(fields)...)
}
func (zl *ZapLogger) Warn(message string, fields map[string]interface{}) {
	zl.logger.Error(message, convertToZapFields(fields)...)
}
func (zl *ZapLogger) Debug(message string, fields map[string]interface{}) {
	zl.logger.Error(message, convertToZapFields(fields)...)
}
func convertToZapFields(fields map[string]interface{}) []zap.Field {
	var zapFields []zap.Field
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return zapFields
}
