package server

import (
	"blog/api/routers"
	"blog/api/validation"
	"blog/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitServer(cfg *config.Config,mysqlCLI *gorm.DB,redisCLI *redis.Client) {
	validation.InitValidations()
	router := routers.InitRouters(cfg.Jwt,mysqlCLI,redisCLI)
	router.Run(":" + string(cfg.Server.Port))
}
