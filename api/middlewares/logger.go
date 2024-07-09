package middlewares

import (
	"blog/config"
	"blog/pkg/logger"
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

type bodyLogWrite struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyLogWrite) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write((b))
}
func (w *bodyLogWrite) WriteString(s string) (int, error) {
	return w.ResponseWriter.WriteString((s))

}


func CustomLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zapLogger := logger.GetZapLoggerInstance(&config.GetConfigInstance().Logger)

		blw := &bodyLogWrite{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
		start := time.Now()
		path := ctx.FullPath()
		raw := ctx.Request.URL.RawQuery
		bodyByte, _ := io.ReadAll(ctx.Request.Body)
		ctx.Request.Body.Close()
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyByte))
		ctx.Writer = blw
		ctx.Next()
		
		param := gin.LogFormatterParams{}
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)
		param.ClientIP = ctx.ClientIP()
		param.Method = ctx.Request.Method
		param.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()
		param.BodySize = ctx.Writer.Size()
		if raw != "" {
			path = path + "?" + raw
		}
		param.Path = path
		keys := map[logger.ExtraKey]interface{}{}
		keys[logger.ClientIp] = param.ClientIP
		keys[logger.Method] = param.Method
		keys[logger.Latency] = param.Latency
		keys[logger.StatusCode] = param.StatusCode
		keys[logger.ErrorMessage] = param.ErrorMessage
		keys[logger.BodySize] = param.BodySize
		keys[logger.RequestBody] = string(bodyByte)
		keys[logger.ResponseBody] = blw.body.String()
		keys[logger.Path] = param.Path

		zapLogger.Info(logger.RequestResponse, logger.Api, "", keys)
	}
}
