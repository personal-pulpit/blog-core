package common

import (
	"blog/api/helpers"
	"blog/pkg/data/models"
	"blog/pkg/data/repo"

	"strconv"

	"github.com/gin-gonic/gin"
)

func IsAdmin(ctx *gin.Context) bool {
	id := helpers.GetIdFromToken(ctx)
	user, err := GetUserFromRedisById(id)
	if err != nil {
		panic(err)
	}
	sRole := user["role"]
	role,_:=strconv.Atoi(sRole)
	return uint(role) == uint(models.AdminRole)
}
func GetUserFromRedisById(id string) (map[string]string, error) {
	ur := repo.NewUserRepo()
	return ur.GetById(id)
}
func GetUserStatus(ctx *gin.Context) bool {
	is_logged := ctx.GetBool("is_logged")
	return is_logged
}