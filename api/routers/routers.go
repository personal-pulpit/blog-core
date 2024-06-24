package routers

import (
	"blog/api/handlers"
	"blog/api/helpers/auth_helper"
	"blog/api/middlewares"
	"blog/config"
	mysql_repository "blog/database/mysql/repo"
	redis_repository "blog/database/redis/repo"
	"blog/internal/service/authentication"
	"blog/internal/service/user"
	"blog/pkg/auth_manager"
	"blog/utils/hash"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var authMiddleware middlewares.UserAuthMiddleware
var authManager auth_manager.AuthManager

func InitRouters(mysqlCLI *gorm.DB, redisCLI *redis.Client) *gin.Engine {
	authManager = auth_manager.NewAuthManager(redisCLI, auth_manager.AuthManagerOpts{PrivateKey: config.Cfg.Jwt.Secret})
	authHelper := auth_helper.NewAuthHeaderHelper()
	authMiddleware = *middlewares.NewUserAuthMiddelware(authManager, authHelper)
	
	
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.CustomLogger())
	r.Use(middlewares.LimitByRequest())


	v1 := r.Group("/api/v1", authMiddleware.SetUserStatus())
	{
		praseRouters(v1.Group(""), mysqlCLI, redisCLI)
		praseRouters(v1.Group("/user"), mysqlCLI, redisCLI)
		praseRouters(v1.Group("/article"), mysqlCLI, redisCLI)
	}

	return r
}
func praseRouters(r *gin.RouterGroup, mysqlCLI *gorm.DB, redisCLI *redis.Client) {

	switch r.BasePath() {
	case "":
		{
			authHelper := auth_helper.NewAuthHeaderHelper()
			h := &handlers.Main{
				AuthHelper: authHelper,
			}
			r.GET("",authMiddleware.SetUserStatus(),h.Main)
		}
	case "/api/v1/user":
		{
			authMysqlRepo := mysql_repository.NewAuthMysqlRepository(mysqlCLI)
			userMysqlRepo := mysql_repository.NewUserMysqlRepository(mysqlCLI)
			hashe_manager := hash.NewHashManager(hash.DefaultHashParams)
			u := &handlers.User{
				AuthService: authentication.NewAuthenticateService(
					authMysqlRepo,
					userMysqlRepo,
					authManager,
					hashe_manager,
				),
				UserService: user.NewUserService(
					userMysqlRepo,
					authMysqlRepo,
				),
			}
			r.GET("/:ID", u.GetByID)
			// r.GET("/logout", middlewares.EnsureLoggedIn(), u.Logout)
			r.POST("", authMiddleware.EnsureNotLoggedIn(), u.Create)
			r.POST("/login", authMiddleware.EnsureNotLoggedIn(), u.Verify)
			r.PATCH("", authMiddleware.EnsureLoggedIn(), u.UpdateByID)
			r.DELETE("/:ID", authMiddleware.EnsureAdmin(), u.DeleteByID)
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
			r.POST("", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), p.Create)
			r.PATCH("", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), p.UpdateByID)
			r.DELETE("/:ID", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), p.DeleteByID)
		}
	}

}
