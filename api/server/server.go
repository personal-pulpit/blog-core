package server

import (
	"blog/api/routers"
	"blog/api/validation"
	"blog/config"
	"blog/pkg/logger"
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitServer(cfg *config.Config, PostgresCLI *gorm.DB, redisCLI *redis.Client,logger logger.Logger)error{
	err := validation.InitValidations()
	if err != nil{
		return err
	}

	router := routers.InitRouters(cfg.Jwt, PostgresCLI, redisCLI,logger)

	return router.Run(fmt.Sprintf(":%d", cfg.Server.Port))
}
