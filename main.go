package main

import (
	"blog/api/server"
	"blog/config"
	postgres "blog/database/postgres"
	redis "blog/database/redis"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/comment"
	"blog/internal/service/user"
	email "blog/pkg/email_manager"
	"blog/pkg/logger"
	"blog/utils/hash"

	postgres_repository "blog/database/postgres/repo"
	_ "blog/docs"
	"blog/pkg/auth_manager"
)

func checkError(loggerInstance logger.Logger, msg string, err error) {
	if err != nil {
		loggerInstance.Fatal(msg, "Error", err)
	}
}

func main() {
	config := config.GetConfigInstance()

	logger := logger.GetZapLoggerInstance()

	postgresCLI, err := postgres.GetPostgresqlDB(&config.Postgres)
	checkError(logger, "Postgres Database", err)

	redisCLI, err := redis.GetRedisDB(&config.Redis)
	checkError(logger, "Redis", err)

	defer redis.CloseRedis()

	articlePostgresRepo := postgres_repository.NewArticlePostgresRepo(postgresCLI)
	authPostgresRepo := postgres_repository.NewAuthPostgresRepository(postgresCLI)
	userPostgresRepo := postgres_repository.NewUserPostgresRepository(postgresCLI)
	commentPostgresRepo := postgres_repository.NewCommentPostgresRepository(postgresCLI)

	hashManager := hash.NewHashManager(hash.DefaultHashParams)
	authManager := auth_manager.NewAuthManger(redisCLI, config.Jwt)

	emailService := email.NewEmailService(&config.Email)
	authService := authentication.NewAuthenticateService(authPostgresRepo, userPostgresRepo, authManager, hashManager, emailService)
	userService := user.NewUserService(userPostgresRepo, authPostgresRepo)
	articleService := article.NewArticleService(articlePostgresRepo)
	commentService := comment.NewCommentService(commentPostgresRepo)

	err = server.InitServer(config.Server.Port, authManager, authService, userService, articleService,commentService, logger)
	checkError(logger, "Server", err)
}
