package routers

import (
	"blog/api/handlers"
	"blog/api/helpers/auth_helper"
	"blog/api/middlewares"

	auth_middlewares "blog/api/middlewares/auth_middlewares"
	"blog/config"
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/repository"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/user"
	"blog/pkg/auth_manager"
	email "blog/pkg/email_manager"
	"blog/pkg/logger"
	"blog/utils/hash"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func init() {
	emailConfigs := config.GetConfigInstance().Email
	emailService = email.NewEmailService(&emailConfigs)
}

var (
	emailService email.EmailService

	authPostgresRepo    repository.AuthPostgresRepository
	userPostgresRepo    repository.UserPostgresRepository
	articlePostgresRepo repository.ArticlePostgresRepository

	authMiddleware *auth_middlewares.UserAuthMiddleware

	authManager auth_manager.AuthManager
	hashManager *hash.HashManager

	authHelper auth_helper.AuthHeaderHelper
)

func InitRouters(jwtCfg config.Jwt, postgresCLI *gorm.DB, redisCLI *redis.Client,logger logger.Logger) *gin.Engine {
	articlePostgresRepo = postgres_repository.NewArticlePostgresRepo(postgresCLI)
	authPostgresRepo = postgres_repository.NewAuthPostgresRepository(postgresCLI)
	userPostgresRepo = postgres_repository.NewUserPostgresRepository(postgresCLI)

	hashManager = hash.NewHashManager(hash.DefaultHashParams)
	authManager = auth_manager.NewAuthManager(redisCLI, auth_manager.AuthManagerOpts{PrivateKey: jwtCfg.Secret})

	authHelper = auth_helper.NewAuthHeaderHelper()

	authMiddleware = auth_middlewares.NewUserAuthMiddelware(authManager, authHelper)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.CustomLogger())
	r.Use(middlewares.LimitByRequest())

	v1 := r.Group("/api/v1", authMiddleware.SetUserStatus())
	{
		praseRouters(v1.Group(""))
		praseRouters(v1.Group("/auth"))
		praseRouters(v1.Group("/user"))
		praseRouters(v1.Group("/article"))
	}

	return r
}
func praseRouters(r *gin.RouterGroup) {

	switch r.BasePath() {
	case "/api/v1":
		{
			mainHandler := &handlers.Main{
				UserService: user.NewUserService(
					userPostgresRepo,
					authPostgresRepo,
				),
			}

			r.GET("", authMiddleware.SetUserStatus(), mainHandler.Main)
		}

	case "/api/v1/auth":
		{
			authHandler := &handlers.AuthHandler{
				AuthService: authentication.NewAuthenticateService(
					authPostgresRepo,
					userPostgresRepo,
					authManager,
					hashManager,
					emailService,
				),
			}

			r.POST("/register", authMiddleware.EnsureNotLoggedIn(), authHandler.Register)
			r.POST("/verifyEmail", authMiddleware.EnsureNotLoggedIn(), authHandler.VerifyEmail)
			r.POST("/login", authMiddleware.EnsureNotLoggedIn(), authHandler.Login)
			r.GET("/logout", authMiddleware.EnsureLoggedIn(), authMiddleware.Logout(), authHandler.Logout)
		}
	case "/api/v1/user":
		{
			userHandler := &handlers.UserHandler{
				UserService: user.NewUserService(
					userPostgresRepo,
					authPostgresRepo,
				),
			}

			r.GET("/:id", userHandler.GetProfile)
			r.PATCH("/update", authMiddleware.EnsureLoggedIn(), userHandler.UpdateProfile)
			r.DELETE("/delete", authMiddleware.EnsureLoggedIn(), userHandler.DeleteAccount)
		}
	case "/api/v1/article":
		{
			articleHandler := &handlers.Article{
				ArticleService: article.NewAricleService(
					articlePostgresRepo,
				),
				UserService: user.NewUserService(
					userPostgresRepo,
					authPostgresRepo,
				),
			}

			r.GET("", articleHandler.GetAll)
			r.GET("/:id", articleHandler.GetById)
			r.POST("", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.Create)
			r.PATCH("/:id", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.UpdateById)
			r.DELETE("/:id", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.DeleteById)
		}
	}
}
