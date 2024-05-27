package common

import (
	"blog/pkg/data/models"
	"blog/pkg/service"
	"blog/utils"

	"github.com/gin-gonic/gin"
)

func IsAdmin(ctx *gin.Context) bool {
	id := utils.GetIdFromToken(ctx)
	user, err := service.GetUserByIdRedis(id)
	if err != nil {
		panic(err)
	}
	return user.Role == uint(models.AdminRole)
}
