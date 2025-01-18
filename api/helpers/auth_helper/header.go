package auth_helper

import (
	"github.com/gin-gonic/gin"
)

const AccessTokenHeader = "X-Access-Token"
const RefreshTokenHeader = "X-Refresh-Token"

func GetHeader(ctx *gin.Context, name string) (string, error) {
	token := ctx.GetHeader(name)
	if token == "" {
		return "", ErrTokenUndefined
	}
	return token, nil
}

func DeleteHeader(ctx *gin.Context, name string) {
	ctx.Request.Header.Del(name)
}
