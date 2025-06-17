package server

import (
	"blog/api/middlewares/auth_middlewares"
	"blog/api/routers"
	"blog/api/validation"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/bookmark"
	"blog/internal/service/category"
	"blog/internal/service/comment"
	"blog/internal/service/like"
	"blog/internal/service/user"
	"blog/pkg/auth_manager"
	"blog/pkg/logger"
	"fmt"
)

type ServerDependencies struct {
	AuthManager     auth_manager.AuthManager
	AuthService     authentication.AuthService
	UserService     user.UserService
	ArticleService  article.ArticleService
	CommentService  comment.CommentService
	CategoryService category.CategoryService
	LikeService     like.LikeService
	BookmarkService bookmark.BookmarkService
	Logger          logger.Logger
}

func NewServerDependencies(
	authManager auth_manager.AuthManager,
	authService authentication.AuthService,
	userService user.UserService,
	articleService article.ArticleService,
	commentService comment.CommentService,
	categoryService category.CategoryService,
	likeService like.LikeService,
	bookmarkService bookmark.BookmarkService,
	logger logger.Logger,
) ServerDependencies {
	return ServerDependencies{
		AuthManager:     authManager,
		AuthService:     authService,
		UserService:     userService,
		ArticleService:  articleService,
		CommentService:  commentService,
		CategoryService: categoryService,
		LikeService:     likeService,
		BookmarkService: bookmarkService,
		Logger:          logger,
	}
}

func InitServer(port int, deps ServerDependencies) error {
	err := validation.InitValidations()
	if err != nil {
		return err
	}

	routerDeps := routers.NewRouterDeps(
		deps.AuthService,
		deps.UserService,
		deps.ArticleService,
		deps.CommentService,
		deps.CategoryService,
		deps.LikeService,
		deps.BookmarkService,
		auth_middlewares.NewUserAuthMiddleware(deps.AuthManager),
		deps.Logger,
	)

	router := routers.InitRouters(
		routerDeps,
	)

	return router.Run(fmt.Sprintf(":%d", port))
}
