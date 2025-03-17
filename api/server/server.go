package server

import (
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

func InitServer(port int, authManager auth_manager.AuthManager,authService authentication.AuthService, userService user.UserService, articleService article.ArticleService,commentService comment.CommentService, logger logger.Logger) error {
	err := validation.InitValidations()
	if err != nil {
		return err
	}

	router := routers.InitRouters(authManager,authService, userService, articleService,commentService,logger)

	return router.Run(fmt.Sprintf(":%d", port))
}
