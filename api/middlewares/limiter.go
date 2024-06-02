package middlewares

import (
	"blog/api/helpers"
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
)

func LimitByRequest() gin.HandlerFunc {
	limiter := tollbooth.NewLimiter(1, nil)
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(limiter, ctx.Writer, ctx.Request)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, helpers.NewErrorHtppResponse(
				http.StatusTooManyRequests, "to many requests", err,
			))
		} else {
			ctx.Next()
		}
	}
}
