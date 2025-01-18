package routers

import (
	"blog/api/handlers"
	"blog/api/middlewares"

	auth_middlewares "blog/api/middlewares/auth_middlewares"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/user"
	"blog/pkg/auth_manager"
	"blog/pkg/logger"

	"github.com/gin-gonic/gin"
)

var (
	authService    authentication.AuthService
	userService    user.UserService
	articleService article.ArticleService

	authMiddleware *auth_middlewares.UserAuthMiddleware
)

func InitRouters(authManager auth_manager.AuthManager,authService authentication.AuthService, userService user.UserService, articleService article.ArticleService, logger logger.Logger) *gin.Engine {
	authService = authService
	userService = userService
	articleService = articleService

	authMiddleware = auth_middlewares.NewUserAuthMiddleware(authManager)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.CustomLogger())
	r.Use(middlewares.LimitByRequest())

	v1 := r.Group("/api/v1", authMiddleware.SetUserStatus())
	{
		praseRouters(v1.Group("/auth"))
		praseRouters(v1.Group("/user"))
		praseRouters(v1.Group("/article"))
	}

	return r
}
func praseRouters(r *gin.RouterGroup) {

	switch r.BasePath() {

	case "/api/v1/auth":
		{
			authHandler := &handlers.AuthHandler{AuthService: authService}

			r.POST("/register", authMiddleware.EnsureNotLoggedIn(), authHandler.Register)
			r.POST("/verifyEmail", authMiddleware.EnsureNotLoggedIn(), authHandler.VerifyEmail)
			r.POST("/login", authMiddleware.EnsureNotLoggedIn(), authHandler.Login)
			r.GET("/logout", authMiddleware.EnsureLoggedIn(), authMiddleware.Logout(), authHandler.Logout)
		}
	case "/api/v1/user":
		{
			userHandler := &handlers.UserHandler{
				UserService: userService,
			}

			r.GET("/:id", userHandler.GetProfile)
			r.PATCH("/update", authMiddleware.EnsureLoggedIn(), userHandler.UpdateProfile)
			r.DELETE("/delete", authMiddleware.EnsureLoggedIn(), userHandler.DeleteAccount)
		}
	case "/api/v1/article":
		{
			articleHandler := &handlers.Article{
				ArticleService: articleService,
				UserService:    userService,
			}

			r.GET("", articleHandler.GetAll)
			r.GET("/:id", articleHandler.GetById)
			r.POST("", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.Create)
			r.PATCH("/:id", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.UpdateById)
			r.DELETE("/:id", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.DeleteById)
		}
	}
}
