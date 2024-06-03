package routers

import (
	"blog/api/handlers"
	"blog/api/middlewares"
	"blog/pkg/data/repo"

	"github.com/gin-gonic/gin"
)

func InitRouters(r *gin.Engine) {
	r.GET("/", middlewares.SetUserStatus(), handlers.Base)
	v1 := r.Group("/api/v1", middlewares.SetUserStatus())
	{
		praseRouters(v1.Group("/user"))
		praseRouters(v1.Group("/article"))
	}
}
func praseRouters(r *gin.RouterGroup) {
	switch r.BasePath() {
	case "/api/v1/user":
		{
			u := &handlers.User{
				UserRepo: repo.NewUserDB(),
			}
			r.GET("", u.GetAll)
			r.GET("/:id", u.GetById)
			r.GET("/logout", middlewares.EnsureLoggedIn(), u.Logout)
			r.POST("", middlewares.EnsureNotLoggedIn(), u.Create)
			r.POST("/login", middlewares.EnsureNotLoggedIn(), u.Verify)
			r.PATCH("", middlewares.EnsureLoggedIn(), u.UpdateById)
			r.DELETE("/:id",middlewares.EnsureAdmin(),u.DeleteById)
		}
	case "/api/v1/article":
		{
			p := &handlers.Article{
				ArticleRepo: repo.NewArticleDB(),
				UserRepo: repo.NewUserDB(),
			}
			r.GET("", p.GetAll)
			r.GET("/:id", p.GetById)
			r.POST("", middlewares.EnsureLoggedIn(),middlewares.EnsureAdmin(), p.Create)
			r.PATCH("",middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.UpdateById)
			r.DELETE("/:id",middlewares.EnsureLoggedIn(),middlewares.EnsureAdmin(),p.DeleteById)
		}
	}

}
