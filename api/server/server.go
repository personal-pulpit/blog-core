package server

import (
	"blog/api/routers"
	"blog/api/validation"
	"blog/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitServer(mysqlCLI *gorm.DB,redisCLI *redis.Client) {
	validation.InitValidations()
	router := routers.InitRouters(mysqlCLI,redisCLI)
	router.Run(":" + config.Cfg.Server.Port)
}
