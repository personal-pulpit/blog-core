package auth_helper

import (
	"github.com/gin-gonic/gin"
)
type AuthHeaderHelper interface {
	GetHeader(ctx *gin.Context)(string,error)
	DeleteHeader(ctx *gin.Context)
}
type authHelperManager struct{}
func NewAuthHeaderHelper()AuthHeaderHelper{
	return &authHelperManager{}
}
func (h *authHelperManager) GetHeader(ctx *gin.Context) (string, error) {
	token := ctx.Writer.Header().Get("Authorization",)
	if token == "" {
		return "", ErrTokenUndefined
	}
	return token, nil
}

func (h *authHelperManager) DeleteHeader(ctx *gin.Context) {
	ctx.Writer.Header().Del("Autherization")
}
