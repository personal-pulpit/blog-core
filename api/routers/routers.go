package routers

import (
	"blog/api/handlers"
	"blog/api/middlewares"
	mysql_repository "blog/database/mysql_repo"

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
				UserRepo: mysql_repository.NewUserRepository(),
			}
			r.GET("", u.GetAll)
			r.GET("/:ID", u.GetByID)
			r.GET("/logout", middlewares.EnsureLoggedIn(), u.Logout)
			r.POST("", middlewares.EnsureNotLoggedIn(), u.Create)
			r.POST("/login", middlewares.EnsureNotLoggedIn(), u.Verify)
			r.PATCH("", middlewares.EnsureLoggedIn(), u.UpdateByID)
			r.DELETE("/:ID", middlewares.EnsureAdmin(), u.DeleteByID)
		}
	case "/api/v1/article":
		{
			p := &handlers.Article{
				UserRepo:    mysql_repository.NewUserRepository(),
				ArticleRepo: mysql_repository.NewArticleRepo(),
			}
			r.GET("", p.GetAll)
			r.GET("/:ID", p.GetByID)
			r.POST("", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.Create)
			r.PATCH("", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.UpdateByID)
			r.DELETE("/:ID", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.DeleteByID)
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
