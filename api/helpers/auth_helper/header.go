package auth_helper

import (
	"github.com/gin-gonic/gin"
)
const AccessTokenHeader = "X-Access-Token"
const VerifyEmailTokenHeader = "X-Verify-Email-Token"
const RefreshTokenHeader = "X-Refresh-Token"
type AuthHeaderHelper interface {
	GetHeader(ctx *gin.Context,name string)(string,error)
	DeleteHeader(ctx *gin.Context,name string)
}
type authHelperManager struct{}
func NewAuthHeaderHelper()AuthHeaderHelper{
	return &authHelperManager{}
}
func (h *authHelperManager) GetHeader(ctx *gin.Context,name string) (string, error) {
	token := ctx.GetHeader(name)
	if token == "" {
		return "", ErrTokenUndefined
	}
	return token, nil
}

func (h *authHelperManager) DeleteHeader(ctx *gin.Context,name string) {
	ctx.Request.Header.Del(name)
}
