package routers

import (
	"blog/api/handlers"
	"blog/api/middlewares"
	mysql_repository "blog/database/mysql/repo"
	redis_repository "blog/database/redis/repo"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func InitRouters(mysqlCLI *gorm.DB, redisCLI *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.CustomLogger())
	r.Use(middlewares.LimitByRequest())
	r.GET("/", middlewares.SetUserStatus(), handlers.Index)
	v1 := r.Group("/api/v1", middlewares.SetUserStatus())
	{
		praseRouters(v1.Group("/user"),mysqlCLI,redisCLI)
		praseRouters(v1.Group("/article"),mysqlCLI,redisCLI)
	}
	return r
}
func praseRouters(r *gin.RouterGroup,mysqlCLI *gorm.DB, redisCLI *redis.Client) {
	switch r.BasePath() {
	case "/api/v1/user":
		{
			u := &handlers.User{
				UserMysqlRepo: mysql_repository.NewUserMysqlRepository(mysqlCLI),
				UserRedisRepo: redis_repository.NewUserRedisRepository(redisCLI),
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
				ArticleMysqlRepo: mysql_repository.NewArticleMysqlRepo(mysqlCLI),
				ArticleRedisRepo: redis_repository.NewArticleRedisRepository(redisCLI),
				UserRedisRepo:    redis_repository.NewUserRedisRepository(redisCLI),
			}
			r.GET("", p.GetAll)
			r.GET("/:ID", p.GetByID)
			r.POST("", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.Create)
			r.PATCH("", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.UpdateByID)
			r.DELETE("/:ID", middlewares.EnsureLoggedIn(), middlewares.EnsureAdmin(), p.DeleteByID)
		}
	}

}

