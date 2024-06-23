package helpers

import (
	"blog/internal/model"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var errFailedSettingToken = errors.New("failed setting token")

func SetToken(ctx *gin.Context, ID model.ID) error {
	if ID == "" {
		return errFailedSettingToken
	}
	token, err := CreateToken(ID)
	if err != nil {
		return err
	}
	ctx.SetCookie("token", token, 3600, "/", "localhost", false, true)
	ctx.Set("is_user_logged", true)
	return nil
}
func GetToken(ctx *gin.Context) jwt.MapClaims {
	tokenStr, err := ctx.Cookie("token")
	if err != nil {
		//check
		panic(err)
	}
	claims, err := ParseToken(tokenStr)
	if err != nil {
		//check
		panic(err)
	}
	return claims
}

func DestroyToken(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
}

func GetIdFromToken(ctx *gin.Context) string {
	claims := GetToken(ctx)
	fid, _ := claims["ID"].(float64)
	sid := strconv.FormatFloat(fid, 'f', -1, 64)
	return sid
}
