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
	"blog/internal/service/category"
	"blog/pkg/logger"

	swaggerFiles "github.com/swaggo/files"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterDeps struct {
	authServiceRouter     authentication.AuthService
	userServiceRouter     user.UserService
	articleServiceRouter  article.ArticleService
	commentServiceRouter  comment.CommentService
	categoryServiceRouter category.CategoryService
	authMiddleware        *auth_middlewares.UserAuthMiddleware
	logger                logger.Logger
}

func NewRouterDeps(
	authServiceRouter authentication.AuthService,
	userServiceRouter user.UserService,
	articleServiceRouter article.ArticleService,
	commentServiceRouter comment.CommentService,
	categoryServiceRouter category.CategoryService,
	authMiddleware *auth_middlewares.UserAuthMiddleware,
	logger logger.Logger,
) RouterDeps {
	return RouterDeps{
		authServiceRouter:    authServiceRouter,
		userServiceRouter:    userServiceRouter,
		articleServiceRouter: articleServiceRouter,
		commentServiceRouter: commentServiceRouter,
		categoryServiceRouter: categoryServiceRouter,
		authMiddleware:       authMiddleware,
		logger:               logger,
	}
}

func InitRouters(routerDeps RouterDeps) *gin.Engine {
	authMiddleware := routerDeps.authMiddleware

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
		praseRouters(v1.Group("/auth"), routerDeps)
		praseRouters(v1.Group("/users"), routerDeps)
		praseRouters(v1.Group("/articles"), routerDeps)
		praseRouters(v1.Group("/comments"), routerDeps)
		praseRouters(v1.Group("/categories"), routerDeps)

	}

	return r
}
func praseRouters(r *gin.RouterGroup, routerDeps RouterDeps) {

	switch r.BasePath() {

	case "/api/v1/auth":
		{
			authHandler := handlers.NewAuthHandler(routerDeps.authServiceRouter)

			r.POST("/register", routerDeps.authMiddleware.EnsureNotLoggedIn(), authHandler.Register)
			r.POST("/verify-email", routerDeps.authMiddleware.EnsureNotLoggedIn(), authHandler.VerifyEmail)
			r.POST("/login", routerDeps.authMiddleware.EnsureNotLoggedIn(), authHandler.Login)
			r.POST("/logout", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.Logout(), authHandler.Logout)
			r.GET("/authenticate", authHandler.Authenticate)
			r.POST("/refresh-token", authHandler.RefreshToken)
			r.POST("/change-password", routerDeps.authMiddleware.EnsureLoggedIn(), authHandler.ChangePassword)
			r.POST("/reset-password/request", routerDeps.authMiddleware.EnsureNotLoggedIn(), authHandler.SendResetPasswordVerification)
			r.POST("/reset-password/submit", routerDeps.authMiddleware.EnsureNotLoggedIn(), authHandler.SubmitResetPassword)
		}

	case "/api/v1/users":
		{
			userHandler := &handlers.UserHandler{
				UserService: routerDeps.userServiceRouter,
			}

			r.GET("/me", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.GetCurrentUser)
			r.POST("/add-profile-picture", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.AddProfileImage)
			r.GET("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.GetUser)
			r.GET("/profile-picture-url/:id", userHandler.GetProfileImageURL)
			r.PATCH("/update", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.UpdateProfile)
			r.DELETE("/delete", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.DeleteAccount)
			r.DELETE("/delete-profile-picture", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.DeleteProfileImage)
		}
	case "/api/v1/articles":
		{
			articleHandler := handlers.NewArticleHandler(routerDeps.articleServiceRouter)

			r.GET("", articleHandler.GetAll)
			r.GET("/:id", articleHandler.GetByID)
			r.GET("/search", articleHandler.Search)
			r.POST("/create", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), articleHandler.Create)
			r.PATCH("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), articleHandler.UpdateByID)
			r.DELETE("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), articleHandler.DeleteByID)
		}
	case "/api/v1/comments":
		{
			commentHandler := handlers.NewCommentHandler(routerDeps.commentServiceRouter)

			r.POST("/create", routerDeps.authMiddleware.EnsureLoggedIn(), commentHandler.Create)
			r.PATCH("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), commentHandler.UpdateByID)
			r.DELETE("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), commentHandler.DeleteByID)
		}
	case "/api/v1/categories":
		{
			categoryHandler := handlers.NewCategoryHandler(routerDeps.categoryServiceRouter)

			r.GET("", categoryHandler.GetAll)
			r.GET("/:id", categoryHandler.GetByID)
			r.POST("/create", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), categoryHandler.Create)
			r.DELETE("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), categoryHandler.DeleteByID)
		}
	}
}
