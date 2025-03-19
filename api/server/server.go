package server

import (
	"blog/api/middlewares/auth_middlewares"
	"blog/api/routers"
	"blog/api/validation"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/comment"
	"blog/internal/service/user"
	"blog/pkg/auth_manager"
	"blog/pkg/logger"
	"fmt"
)

type ServerDependencies struct {
	AuthManager    auth_manager.AuthManager
	AuthService    authentication.AuthService
	UserService    user.UserService
	ArticleService article.ArticleService
	CommentService comment.CommentService
	Logger         logger.Logger
}

func NewServerDependencies(
	authManager auth_manager.AuthManager,
	authService authentication.AuthService,
	userService user.UserService,
	articleService article.ArticleService,
	commentService comment.CommentService,
	logger logger.Logger,
) ServerDependencies {
	return ServerDependencies{
		AuthManager:    authManager,
		AuthService:    authService,
		UserService:    userService,
		ArticleService: articleService,
		CommentService: commentService,
		Logger:         logger,
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
		auth_middlewares.NewUserAuthMiddleware(deps.AuthManager),
		deps.Logger,
	)

	router := routers.InitRouters(
		routerDeps,
	)

	return router.Run(fmt.Sprintf(":%d", port))
}
