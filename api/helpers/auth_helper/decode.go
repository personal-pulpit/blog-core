package auth_helper

import (
	"blog/config"
	database "blog/database/redis"
	"blog/internal/model"
	"blog/pkg/auth_manager"
)
var authManager auth_manager.AuthManager
func init(){
	redisClient := database.GetRedisDB() 
	authManager = auth_manager.NewAuthManager(redisClient,auth_manager.AuthManagerOpts{
		PrivateKey: config.Cfg.Jwt.Secret,
	})
}
func GetIdByToken(token string,tokenType auth_manager.TokenType)(model.ID,error){
	tokenClaims,err := authManager.DecodeToken(token,tokenType)
	if err != nil{
		return "",err
	}
	return tokenClaims.ID,err
}