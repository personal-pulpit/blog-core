package middlewares

import (
	"blog/utils"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)
var(
	ErrYouAreUnAuthorized = errors.New("you are unauthorized")
)
func SetUserStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := ctx.Cookie("token")
		if err != nil || token == "" {
			ctx.Set("is_logged", false)
			ctx.Next()
		} else {
			ctx.Set("is_logged", true)
			ctx.Set("is_admin", utils.IsAdmin(ctx))
			ctx.Next()
		}
	}
}
func EnsureLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		is_logged := ctx.GetBool("is_logged")
		if is_logged {
			ctx.Next()
		} else {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			utils.NewErrorHtppResponse(http.StatusUnauthorized,
				"sometime went wrong",ErrYouAreUnAuthorized))
			return
		}
	}
}
func EnsureNotLoggedIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		is_logged := ctx.GetBool("is_logged")
		if !is_logged {
			ctx.Next()
		} else {
			utils.DestroyToken(ctx)
			ctx.Next()
		}
	}
}
func EnsureAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !utils.IsAdmin(ctx) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized,
				utils.NewErrorHtppResponse(http.StatusUnauthorized,
					"sometime went wrong",ErrYouAreUnAuthorized))
				return
		} else {
			ctx.Next()
		}
	}
}
