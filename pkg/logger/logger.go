package logger

type Logger interface {
    Info(message string, fields map[string]interface{})
    Error(message string, fields map[string]interface{})
    Warn(message string, fields map[string]interface{})
    Debug(message string, fields map[string]interface{})
}