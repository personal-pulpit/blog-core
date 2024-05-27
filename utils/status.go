package utils

import (
	"github.com/gin-gonic/gin"
)

func GetUserStatus(ctx *gin.Context) bool {
	is_logged := ctx.GetBool("is_logged")
	return is_logged
}
