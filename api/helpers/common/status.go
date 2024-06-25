package common

import (
	"blog/internal/model"

	"github.com/gin-gonic/gin"
)

func IsAdmin(role model.Role) bool {
	return role == model.AdminRole
}
func GetUserStatus(ctx *gin.Context) bool {
	is_logged := ctx.GetBool("is_logged")
	return is_logged
}
