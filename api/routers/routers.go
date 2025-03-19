package routers

import (
	"blog/api/handlers"
	"blog/api/middlewares"
	"blog/config"

	auth_middlewares "blog/api/middlewares/auth_middlewares"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/comment"
	"blog/internal/service/user"
	"blog/pkg/auth_manager"
	"blog/pkg/logger"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	authServiceRouter    authentication.AuthService
	userServiceRouter    user.UserService
	articleServiceRouter article.ArticleService
	commentServiceRouter comment.CommentService


	authMiddleware *auth_middlewares.UserAuthMiddleware
)

func InitRouters(authManager auth_manager.AuthManager, authService authentication.AuthService, userService user.UserService, articleService article.ArticleService,commentService comment.CommentService, logger logger.Logger) *gin.Engine {
	authServiceRouter = authService
	userServiceRouter = userService
	articleServiceRouter = articleService
	commentServiceRouter = commentService

	authMiddleware = auth_middlewares.NewUserAuthMiddleware(authManager)

	if config.GetEnv() == config.Production {
		gin.SetMode(gin.ReleaseMode)
	}
	
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.LimitByRequest(1))

	//swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1", authMiddleware.SetUserStatus())
	{
		praseRouters(v1.Group("/auth"))
		praseRouters(v1.Group("/users"))
		praseRouters(v1.Group("/articles"))
		praseRouters(v1.Group("/comments"))
	}

	return r
}
func praseRouters(r *gin.RouterGroup) {

	switch r.BasePath() {

	case "/api/v1/auth":
		{
			authHandler := handlers.NewAuthHandler(authServiceRouter)

			r.POST("/register", authMiddleware.EnsureNotLoggedIn(), authHandler.Register)
			r.POST("/verify-email", authMiddleware.EnsureNotLoggedIn(), authHandler.VerifyEmail)
			r.POST("/login", authMiddleware.EnsureNotLoggedIn(), authHandler.Login)
			r.POST("/logout", authMiddleware.EnsureLoggedIn(), authMiddleware.Logout(), authHandler.Logout)
			r.GET("/authenticate", authHandler.Authenticate)
			r.POST("/refresh-token", authHandler.RefreshToken)
			r.POST("/change-password", authMiddleware.EnsureLoggedIn(), authHandler.ChangePassword)
			r.POST("/reset-password/request", authMiddleware.EnsureNotLoggedIn(), authHandler.SendResetPasswordVerification)
			r.POST("/reset-password/submit", authMiddleware.EnsureNotLoggedIn(), authHandler.SubmitResetPassword)
		}

	case "/api/v1/users":
		{
			userHandler := &handlers.UserHandler{
				UserService: userServiceRouter,
			}

			r.GET("/me", authMiddleware.EnsureLoggedIn(), userHandler.GetCurrentUser)
			r.GET("/:id", authMiddleware.EnsureLoggedIn(), userHandler.GetUser)
			r.PATCH("/update", authMiddleware.EnsureLoggedIn(), userHandler.UpdateProfile)
			r.DELETE("/delete", authMiddleware.EnsureLoggedIn(), userHandler.DeleteAccount)
		}
	case "/api/v1/articles":
		{
			articleHandler := handlers.NewArticleHandler(articleServiceRouter)

			r.GET("", articleHandler.GetAll)
			r.GET("/:id", articleHandler.GetByID)
			r.GET("/search", articleHandler.Search)
			r.POST("/create", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.Create)
			r.PATCH("/:id", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.UpdateByID)
			r.DELETE("/:id", authMiddleware.EnsureLoggedIn(), authMiddleware.EnsureAdmin(), articleHandler.DeleteByID)
		}
	case "/api/v1/comments":
	{
		commentHandler := handlers.NewCommentHandler(commentServiceRouter)

		r.POST("/create", authMiddleware.EnsureLoggedIn(), commentHandler.Create)
		r.PATCH("/:id", authMiddleware.EnsureLoggedIn(), commentHandler.UpdateByID)
		r.DELETE("/:id", authMiddleware.EnsureLoggedIn(), commentHandler.DeleteByID)
	}
}
}
