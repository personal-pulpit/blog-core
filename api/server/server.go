package server

import (
	"blog/api/routers"

	"github.com/gin-gonic/gin"
)

func InitServer(){
	r := gin.New()
	r.Use(gin.Logger(),gin.Recovery())
	routers.InitRouters(r)
	r.Run(":8080")
}