package routers

import (
	"blog/api/handlers"
	"blog/api/middlewares"

	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.CustomLogger())
	r.Use(middlewares.LimitByRequest())
	r.GET("/", middlewares.SetUserStatus(), handlers.Index)
	v1 := r.Group("/api/v1", middlewares.SetUserStatus())
	{
		praseRouters(v1.Group("/user"))
		praseRouters(v1.Group("/article"))
	}
	return r
}
func praseRouters(r *gin.RouterGroup) {
	switch r.BasePath() {
	case "/api/v1/user":
		{
			u := &handlers.User{
				
			}
			r.GET("", u.GetAll)
			r.GET("/:id", u.GetById)
			r.GET("/logout", middlewares.EnsureLoggedIn(), u.Logout)
			r.POST("", middlewares.EnsureNotLoggedIn(), u.Create)
			r.POST("/login", middlewares.EnsureNotLoggedIn(), u.Verify)
			r.PATCH("", middlewares.EnsureLoggedIn(), u.UpdateById)
			r.DELETE("/:id", middlewares.EnsureAdmin(), u.DeleteById)
		}
	case "/api/v1/article":
		{
			p := &handlers.Article{}
			r.GET("", p.GetAll)
			r.GET("/:id", p.GetById)
			r.POST("", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.Create)
			r.PATCH("", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.UpdateById)
			r.DELETE("/:id", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.DeleteById)
		}
	}

}
func InitRoutersForTest() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.LimitByRequest())
	r.GET("/", middlewares.SetUserStatus(), handlers.Index)
	v1 := r.Group("/api/v1", middlewares.SetUserStatus())
	{
		praseRouters(v1.Group("/user"))
		praseRouters(v1.Group("/article"))
	}
	return r
}
