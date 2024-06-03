package server

import (
	"blog/api/middlewares"
	"blog/api/routers"
	"blog/api/validation"
	"blog/config"

	"github.com/gin-gonic/gin"
)

func InitServer() {
	r := gin.New()
	validation.InitValidations()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.CustomLogger())
	r.Use(middlewares.LimitByRequest())
	routers.InitRouters(r)
	r.Run(":" + config.Cfg.Server.Port)
}
