package routers

import (
	"blog/api/handlers"
	"blog/api/middlewares"
	"blog/config"

	auth_middlewares "blog/api/middlewares/auth_middlewares"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/bookmark"
	"blog/internal/service/category"
	"blog/internal/service/comment"
	"blog/internal/service/like"
	"blog/internal/service/user"
	"blog/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouterDeps struct {
	authServiceRouter     authentication.AuthService
	userServiceRouter     user.UserService
	articleServiceRouter  article.ArticleService
	commentServiceRouter  comment.CommentService
	categoryServiceRouter category.CategoryService
	likeServiceRouter     like.LikeService
	bookmarkServiceRouter bookmark.BookmarkService
	authMiddleware        *auth_middlewares.UserAuthMiddleware
	logger                logger.Logger
}

func NewRouterDeps(
	authServiceRouter authentication.AuthService,
	userServiceRouter user.UserService,
	articleServiceRouter article.ArticleService,
	commentServiceRouter comment.CommentService,
	categoryServiceRouter category.CategoryService,
	likeServiceRouter like.LikeService,
	bookmarkServiceRouter bookmark.BookmarkService,
	authMiddleware *auth_middlewares.UserAuthMiddleware,
	logger logger.Logger,
) RouterDeps {
	return RouterDeps{
		authServiceRouter:     authServiceRouter,
		userServiceRouter:     userServiceRouter,
		articleServiceRouter:  articleServiceRouter,
		commentServiceRouter:  commentServiceRouter,
		categoryServiceRouter: categoryServiceRouter,
		likeServiceRouter:     likeServiceRouter,
		bookmarkServiceRouter: bookmarkServiceRouter,
		authMiddleware:        authMiddleware,
		logger:                logger,
	}
}

func InitRouters(routerDeps RouterDeps) *gin.Engine {
	authMiddleware := routerDeps.authMiddleware

	if config.GetEnv() == config.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middlewares.LimitByRequest(4))

	configInstance := config.GetConfigInstance()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{configInstance.Server.CORS}                // Allow your frontend origin
	config.AllowMethods = []string{"POST", "OPTIONS", "GET", "PUT", "DELETE"} // Allow required methods
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"} // Allow required headers
	config.AllowCredentials = true

	r.Use(cors.New(config))
	//swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1", authMiddleware.SetUserStatus())

	{
		praseAuthRouters(v1.Group("/auth"), routerDeps)
		praseUserRouters(v1.Group("/users"), routerDeps)
		praseArticleRouters(v1.Group("/articles"), routerDeps)
		praseCommentRouters(v1.Group("/comments"), routerDeps)
		praseCategoryRouters(v1.Group("/categories"), routerDeps)
		praseLikeRouters(v1.Group("/likes"), routerDeps)
		praseBookmarkRouters(v1.Group("/bookmarks"), routerDeps)
	}

	return r
}

func praseAuthRouters(r *gin.RouterGroup, routerDeps RouterDeps) {
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

func praseUserRouters(r *gin.RouterGroup, routerDeps RouterDeps) {
	userHandler := &handlers.UserHandler{
		UserService: routerDeps.userServiceRouter,
	}

	r.GET("", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.GetCurrentUser)
	r.POST("/profile-picture", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.AddProfileImage)
	r.GET("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.GetUser)
	r.GET("/:id/profile-picture-url", userHandler.GetProfileImageURL)
	r.PATCH("", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.UpdateProfile)
	r.DELETE("", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.DeleteAccount)
	r.DELETE("/profile-picture", routerDeps.authMiddleware.EnsureLoggedIn(), userHandler.DeleteProfileImage)
}

func praseArticleRouters(r *gin.RouterGroup, routerDeps RouterDeps) {
	articleHandler := handlers.NewArticleHandler(routerDeps.articleServiceRouter)

	r.GET("", articleHandler.GetArticles)
	r.GET("/:id", articleHandler.GetByID)
	r.POST("", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), articleHandler.Create)
	r.PATCH("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), articleHandler.UpdateByID)
	r.DELETE("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), articleHandler.DeleteByID)
}

func praseCommentRouters(r *gin.RouterGroup, routerDeps RouterDeps) {
	commentHandler := handlers.NewCommentHandler(routerDeps.commentServiceRouter)

	r.POST("", routerDeps.authMiddleware.EnsureLoggedIn(), commentHandler.Create)
	r.PATCH("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), commentHandler.UpdateByID)
	r.DELETE("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), commentHandler.DeleteByID)
}

func praseCategoryRouters(r *gin.RouterGroup, routerDeps RouterDeps) {
	categoryHandler := handlers.NewCategoryHandler(routerDeps.categoryServiceRouter)

	r.GET("", categoryHandler.GetAll)
	r.GET("/:id", categoryHandler.GetByID)
	r.POST("", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), categoryHandler.Create)
	r.DELETE("/:id", routerDeps.authMiddleware.EnsureLoggedIn(), routerDeps.authMiddleware.EnsureAdmin(), categoryHandler.DeleteByID)
}

func praseLikeRouters(r *gin.RouterGroup, routerDeps RouterDeps) {
	likeHandler := handlers.NewLikeHandler(routerDeps.likeServiceRouter)

	r.GET("", routerDeps.authMiddleware.EnsureLoggedIn(), likeHandler.GetUserLikes)
	r.POST("", routerDeps.authMiddleware.EnsureLoggedIn(), likeHandler.Create)	
	r.DELETE("/article/:id", routerDeps.authMiddleware.EnsureLoggedIn(), likeHandler.DeleteByID)
}

func praseBookmarkRouters(r *gin.RouterGroup, routerDeps RouterDeps) {
	bookmarkHandler := handlers.NewBookmarkHandler(routerDeps.bookmarkServiceRouter)
	r.GET("", routerDeps.authMiddleware.EnsureLoggedIn(), bookmarkHandler.GetUserBookmarks)
	r.POST("", routerDeps.authMiddleware.EnsureLoggedIn(), bookmarkHandler.Create)
	r.DELETE("/article/:id", routerDeps.authMiddleware.EnsureLoggedIn(), bookmarkHandler.Delete)
}
