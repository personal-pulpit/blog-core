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
	authServiceRouter    authentication.AuthService
	userServiceRouter   user.UserService
	articleServiceRouter article.ArticleService

	authMiddleware *auth_middlewares.UserAuthMiddleware
)

func InitRouters(authManager auth_manager.AuthManager, authService authentication.AuthService, userService user.UserService, articleService article.ArticleService, logger logger.Logger) *gin.Engine {
	authServiceRouter = authService
	userServiceRouter = userService
	articleServiceRouter = articleService

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
			authHandler := handlers.NewAuthHandler(authServiceRouter)

			r.POST("/register", authMiddleware.EnsureNotLoggedIn(), authHandler.Register)
			r.POST("/verifyEmail", authMiddleware.SetUserStatus(), authMiddleware.EnsureNotLoggedIn(), authHandler.VerifyEmail)
			r.POST("/login", authMiddleware.SetUserStatus(), authMiddleware.EnsureNotLoggedIn(), authHandler.Login)
			r.GET("/logout", authMiddleware.SetUserStatus(),authMiddleware.EnsureLoggedIn(), authMiddleware.Logout(), authHandler.Logout)
			r.GET("/authenticate", authHandler.Authenticate)
			r.POST("/refresh-token", authHandler.RefreshToken)
			r.POST("/change-password",  authMiddleware.SetUserStatus(),authMiddleware.EnsureLoggedIn(), authHandler.ChangePassword)
			r.POST("/reset-password/request", authMiddleware.SetUserStatus(), authMiddleware.EnsureNotLoggedIn(), authHandler.SendResetPasswordVerification)
			r.POST("/reset-password/submit", authMiddleware.SetUserStatus(), authMiddleware.EnsureNotLoggedIn(), authHandler.SubmitResetPassword)
		}

	case "/api/v1/user":
		{
			userHandler := &handlers.UserHandler{
				UserService: userServiceRouter,
			}

			r.GET("/me",authMiddleware.SetUserStatus(),authMiddleware.EnsureLoggedIn(), userHandler.GetUser)
			r.GET("/:id",authMiddleware.SetUserStatus(),authMiddleware.EnsureLoggedIn(), userHandler.GetUser)
			r.PATCH("/update",authMiddleware.SetUserStatus(), authMiddleware.EnsureLoggedIn(), userHandler.UpdateProfile)
			r.DELETE("/delete",authMiddleware.SetUserStatus(), authMiddleware.EnsureLoggedIn(), userHandler.DeleteAccount)
		}
	case "/api/v1/article":
		{
			articleHandler := handlers.NewArticleHandler(articleServiceRouter)

			r.GET("",authMiddleware.SetUserStatus(), articleHandler.GetAll)
			r.GET("/:id",authMiddleware.SetUserStatus(), articleHandler.GetByID)
			r.POST("", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.Create)
			r.PATCH("/:id", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.UpdateByID)
			r.DELETE("/:id", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.DeleteByID)
		}
	}
}
