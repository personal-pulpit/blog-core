package server

import (
	"blog/api/routers"
	"blog/api/validation"
	"blog/config"
)

func InitServer() {
	validation.InitValidations()
	router := routers.InitRouters()
	router.Run(":" + config.Cfg.Server.Port)
}
