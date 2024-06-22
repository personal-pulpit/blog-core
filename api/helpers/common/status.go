package common

import (
	"blog/api/helpers"
	"blog/pkg/data/model"
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
	role, _ := strconv.Atoi(sRole)
	return uint(role) == uint(model.AdminRole)
}
func GetUserFromRedisById(id string) (map[string]string, error) {
	ur := repo.NewUserDB()
	return ur.GetById(id)
}
func GetUserStatus(ctx *gin.Context) bool {
	is_logged := ctx.GetBool("is_logged")
	return is_logged
}
