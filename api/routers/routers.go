package routers

import (
	"blog/api/handlers"
	"blog/api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	r.GET("/",middlewares.SetUserStatus(), handlers.Base)
	v1 := r.Group("/api/v1",middlewares.SetUserStatus())
	{
		praseRouters(v1.Group("/user"))
	}
}
func praseRouters(r *gin.RouterGroup) {
	switch r.BasePath() {
	case "/api/v1/user":
		{
			u := handlers.User{}
			r.GET("",u.GetUsers)
			r.GET("/:id", u.Get)
			r.GET("/logout", middlewares.EnsureLoggedIn(),u.Logout)
			r.POST("",middlewares.EnsureNotLoggedIn(),u.Create)
			r.POST("/login",middlewares.EnsureNotLoggedIn(),u.Verify)
			r.PATCH("",middlewares.EnsureLoggedIn(),u.UpdateById)
			r.DELETE("/:id", u.DeleteById)
		}	
	}

}
