package routers

import (
	"blog/api/handlers"

	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	r.GET("/", handlers.Base)
	v1 := r.Group("/api/v1")
	{
		praseRouters(v1.Group("/user"))
	}
}
func praseRouters(r *gin.RouterGroup) {
	switch r.BasePath() {
	case "/api/v1/user":
		{
			u := handlers.User{}
			r.GET("/:id", u.Get)
			r.POST("", u.Create)
			r.PATCH("", u.UpdateById)
			r.DELETE("/:id", u.DeleteById)
		}
	}

}
