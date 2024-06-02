package common

import (
	"blog/pkg/data/models"
	"blog/pkg/data/repo"
	"blog/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func IsAdmin(ctx *gin.Context) bool {
	id := utils.GetIdFromToken(ctx)
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
