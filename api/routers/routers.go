package routers

import (
	"blog/api/handlers"
	"github.com/gin-gonic/gin"
)


func InitRouters(r *gin.Engine){
	r.GET("/",handlers.Base)
	v1 := r.Group("/v1")
	{
		v1.GET("/:id",handlers.Get)
		v1.DELETE("/:id",handlers.Delete)
		v1.POST("",handlers.Post)
		v1.PUT("",handlers.Put)

	}
	
}