package helpers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// time must be second
func NewContextWithTimeout(ctx *gin.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx.Request.Context(), timeout*time.Second)
}
func GetResponse(ctx *gin.Context, baseCode int, ResponseChannel <-chan HttpResponse) {
	ctxWithTimeout, cancel := NewContextWithTimeout(ctx, 5)
	defer cancel()
	select {
	case response := <-ResponseChannel:
		if response.Code == baseCode {
			ctx.JSON(baseCode, response)
			return
		}
		ctx.AbortWithStatusJSON(response.Code, response)
	case <-ctxWithTimeout.Done():
		ctx.AbortWithStatusJSON(http.StatusRequestTimeout, NewHttpResponse(
			http.StatusBadRequest, "timed out", nil),
		)
	}
}
