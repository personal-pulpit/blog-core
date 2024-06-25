package routers

import (
	"blog/api/handlers"
	"blog/api/helpers/auth_helper"
	"blog/api/middlewares"
	"blog/config"
	mysql_repository "blog/database/mysql/repo"
	"blog/internal/repository"
	"blog/internal/service/authentication"
	"blog/internal/service/user"
	"blog/pkg/auth_manager"
	"blog/utils/hash"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	authMysqlRepo  repository.AuthMysqlRepository
	userMysqlRepo  repository.UserMysqlRepository
	authMiddleware *middlewares.UserAuthMiddleware
	authManager    auth_manager.AuthManager
	hashManager    *hash.HashManager
)

func InitRouters(jwtCfg config.Jwt, mysqlCLI *gorm.DB, redisCLI *redis.Client) *gin.Engine {
	authMysqlRepo = mysql_repository.NewAuthMysqlRepository(mysqlCLI)
	userMysqlRepo = mysql_repository.NewUserMysqlRepository(mysqlCLI)
	hashManager = hash.NewHashManager(hash.DefaultHashParams)
	authManager = auth_manager.NewAuthManager(redisCLI, auth_manager.AuthManagerOpts{PrivateKey: jwtCfg.Secret})
	authHelper := auth_helper.NewAuthHeaderHelper()
	authMiddleware = middlewares.NewUserAuthMiddelware(authManager, authHelper)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.CustomLogger())
	r.Use(middlewares.LimitByRequest())

	v1 := r.Group("/api/v1", authMiddleware.SetUserStatus())
	{
		praseRouters(v1.Group(""))
		praseRouters(v1.Group("/user"))
		praseRouters(v1.Group("/article"))
	}

	return r
}
func praseRouters(r *gin.RouterGroup,) {

	switch r.BasePath() {
	case "":
		{
			mainHandler := &handlers.Main{
				UserService: user.NewUserService(
					userMysqlRepo,
					authMysqlRepo,
				),
			}
			r.GET("", authMiddleware.SetUserStatus(), mainHandler.Main)
		}
	case "/api/v1/user":
		{
			userHandler := &handlers.User{
				AuthService: authentication.NewAuthenticateService(
					authMysqlRepo,
					userMysqlRepo,
					authManager,
					hashManager,
				),
				UserService: user.NewUserService(
					userMysqlRepo,
					authMysqlRepo,
				),
			}
			r.GET("/:id", userHandler.GetProfile)
			r.POST("", authMiddleware.EnsureNotLoggedIn(), userHandler.Register)
			r.POST("/login", authMiddleware.EnsureNotLoggedIn(), userHandler.Login)
			r.PATCH("", authMiddleware.EnsureLoggedIn(), userHandler.UpdateProfile)
			r.DELETE("", authMiddleware.EnsureAdmin(), userHandler.DeleteAccount)
			r.GET("/logout", authMiddleware.EnsureLoggedIn(), userHandler.Logout)
		}
	}
}

// case "/api/v1/article":
// 	{
// 		articleHandler := &handlers.Article{
// 			ArticleMysqlRepo: mysql_repository.NewArticleMysqlRepo(mysqlCLI),
// 			ArticleRedisRepo: redis_repository.NewArticleRedisRepository(redisCLI),
// 			UserRedisRepo:    redis_repository.NewUserRedisRepository(redisCLI),
// 		}
// 		r.GET("", articleHandler.GetAll)
// 		r.GET("/:ID", articleHandler.GetByID)
// 		r.POST("", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.Create)
// 		r.PATCH("", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.UpdateByID)
// 		r.DELETE("/:ID", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.DeleteByID)
// 	}
// }

// }
