package utils

import "github.com/gin-gonic/gin"

func IsAdmin(ctx *gin.Context) bool {
	claims := GetToken(ctx)
	role := claims["role"]
	rolef, _ := role.(float64)
	roleu := uint(rolef)
	return roleu == 2
}
func GetUserStatus(ctx *gin.Context) bool {
	is_logged := ctx.GetBool("is_logged")
	return is_logged
}
