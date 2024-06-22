package common

import (
	"blog/api/helpers"
	mysql_repository "blog/database/mysql_repo"
	"blog/internal/model"

	"strconv"

	"github.com/gin-gonic/gin"
)

func IsAdmin(ctx *gin.Context) bool {
	ID := helpers.GetIdFromToken(ctx)
	user, err := GetUserFromRedisById(ID)
	if err != nil {
		panic(err)
	}
	sRole := user["role"]
	role, _ := strconv.Atoi(sRole)
	return uint(role) == uint(model.AdminRole)
}
func GetUserFromRedisById(ID string) (map[string]string, error) {
	u := mysql_repository.NewUserRepository()
	return u.GetByID(ID)
}
func GetUserStatus(ctx *gin.Context) bool {
	is_logged := ctx.GetBool("is_logged")
	return is_logged
}
