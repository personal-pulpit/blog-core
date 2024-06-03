package middlewares

import (
	"blog/pkg/logging"
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
		keys := map[logging.ExtraKey]interface{}{}
		keys[logging.ClientIp] = param.ClientIP
		keys[logging.Method] = param.Method
		keys[logging.Latency] = param.Latency
		keys[logging.StatusCode] = param.StatusCode
		keys[logging.ErrorMessage] = param.ErrorMessage
		keys[logging.BodySize] = param.BodySize
		keys[logging.RequestBody] = string(bodyByte)
		keys[logging.ResponseBody] = blw.body.String()
		logging.MyLogger.Info(logging.RequestResponse, logging.Api, "", keys)
	}
}
