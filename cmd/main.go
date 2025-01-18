package main

import (
	"blog/api/server"
	"blog/config"
	postgres "blog/database/postgres"
	redis "blog/database/redis"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/user"
	email "blog/pkg/email_manager"
	"blog/pkg/logger"
	"blog/utils/hash"

	postgres_repository "blog/database/postgres/repo"
	"blog/pkg/auth_manager"
)

func checkError(loggerInstance logger.Logger, err error) {
	if err != nil {
		loggerInstance.Fatal(logger.General, logger.Startup, err.Error(), map[logger.ExtraKey]interface{}{})
	}
}
func main() {
	config := config.GetConfigInstance()

	logger := logger.GetZapLoggerInstance(&config.Logger)

	postgresCLI, err := postgres.GetPostgresqlDB(&config.Postgres)
	checkError(logger, err)

	redisCLI, err := redis.GetRedisDB(&config.Redis)
	checkError(logger, err)

	defer redis.CloseRedis()

	articlePostgresRepo := postgres_repository.NewArticlePostgresRepo(postgresCLI)
	authPostgresRepo := postgres_repository.NewAuthPostgresRepository(postgresCLI)
	userPostgresRepo := postgres_repository.NewUserPostgresRepository(postgresCLI)

	hashManager := hash.NewHashManager(hash.DefaultHashParams)
	authManager := auth_manager.NewAuthManger(redisCLI, config.Jwt)

	emailService := email.NewEmailService(&config.Email)
	authService := authentication.NewAuthenticateService(authPostgresRepo, userPostgresRepo, authManager, hashManager, emailService)
	userService := user.NewUserService(userPostgresRepo)
	articleService := article.NewArticleService(articlePostgresRepo)

	err = server.InitServer(config.Server.Port, authManager,authService, userService, articleService, logger)
	checkError(logger, err)
}
