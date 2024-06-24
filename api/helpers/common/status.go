package common

import (
	"blog/api/helpers/auth_helper"
	database "blog/database/redis"
	redis_repository "blog/database/redis/repo"
	"blog/pkg/auth_manager"

	"blog/internal/model"

	"strconv"

	"github.com/gin-gonic/gin"
)

func IsAdmin(token string, tokenType auth_manager.TokenType) bool {
	id, _ := auth_helper.GetIdByToken(token, tokenType)
	user, err := GetUserFromRedisByID(id)
	if err != nil {
		panic(err)
	}
	sRole := user["role"]
	role, _ := strconv.Atoi(sRole)
	return uint(role) == uint(model.AdminRole)
}
func GetUserFromRedisByID(id string) (map[string]string, error) {
	redisCLI := database.GetRedisDB()
	u := redis_repository.NewUserRedisRepository(redisCLI)
	return u.GetCacheByID(id)
}
func GetUserStatus(ctx *gin.Context) bool {
	is_logged := ctx.GetBool("is_logged")
	return is_logged
}
