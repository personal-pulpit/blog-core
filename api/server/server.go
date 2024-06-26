package server

import (
	"blog/api/routers"
	"blog/api/validation"
	"blog/config"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitServer(cfg *config.Config, PostgresCLI *gorm.DB, redisCLI *redis.Client) {
	validation.InitValidations()
	router := routers.InitRouters(cfg.Jwt, PostgresCLI, redisCLI)
	router.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
