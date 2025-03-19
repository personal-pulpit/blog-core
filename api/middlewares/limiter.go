package middlewares

import (
	"blog/api/helpers"
	"net/http"

	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
)

func LimitByRequest(limit float64) gin.HandlerFunc {
	limiter := tollbooth.NewLimiter(limit, nil)
	return func(ctx *gin.Context) {
		err := tollbooth.LimitByRequest(limiter, ctx.Writer, ctx.Request)
		if err != nil {
			helpers.RespondWithError(ctx, http.StatusTooManyRequests, err.Error(), nil)
			return
		}
		ctx.Next()
	}
}