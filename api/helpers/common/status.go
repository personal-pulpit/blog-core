package common

import (
	"blog/api/helpers"
	redis_repository "blog/database/redis/repo"

	"blog/internal/model"

	"strconv"

	"github.com/gin-gonic/gin"
)

func IsAdmin(ctx *gin.Context) bool {
	ID := helpers.GetIdFromToken(ctx)
	user, err := GetUserFromRedisByID(ID)
	if err != nil {
		panic(err)
	}
	sRole := user["role"]
	role, _ := strconv.Atoi(sRole)
	return uint(role) == uint(model.AdminRole)
}
func GetUserFromRedisByID(ID string) (map[string]string, error) {
	u := redis_repository.NewUserRedisRepository()
	return u.GetCacheByID(ID)
}
func GetUserStatus(ctx *gin.Context) bool {
	is_logged := ctx.GetBool("is_logged")
	return is_logged
}
