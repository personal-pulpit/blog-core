package main

import (
	"blog/api/server"
	"blog/config"
	postgres "blog/database/postgres"
	redis "blog/database/redis"
	redis_repo "blog/database/redis/repo"
	"blog/internal/service/article"
	"blog/internal/service/authentication"
	"blog/internal/service/bookmark"
	"blog/internal/service/category"
	"blog/internal/service/comment"
	"blog/internal/service/like"
	"blog/internal/service/user"
	email "blog/pkg/email_manager"
	"blog/pkg/logger"
	objStorage "blog/pkg/object_storage"
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

	minioClient, err := objStorage.NewMinioClient(config.MinIO)
	checkError(logger, "MinIO", err)

	articlePostgresRepo := postgres_repository.NewArticlePostgresRepo(postgresCLI)
	authPostgresRepo := postgres_repository.NewAuthPostgresRepository(postgresCLI)
	userPostgresRepo := postgres_repository.NewUserPostgresRepository(postgresCLI)
	userCacheRepo := redis_repo.NewRedisRepository(redisCLI)
	commentPostgresRepo := postgres_repository.NewCommentPostgresRepository(postgresCLI)
	categoryPostgresRepo := postgres_repository.NewCategoryRepository(postgresCLI)
	likePostgresRepo := postgres_repository.NewLikePostgresRepository(postgresCLI)
	bookmarkPostgresRepo := postgres_repository.NewBookmarkPostgresRepository(postgresCLI)
	objStorageRepo := objStorage.NewMinioStorageRepo(minioClient)

	hashManager := hash.NewHashManager(hash.DefaultHashParams)
	authManager := auth_manager.NewAuthManger(redisCLI, config.Jwt)

	emailService := email.NewEmailService(&config.Email)
	authService := authentication.NewAuthenticateService(authPostgresRepo, userPostgresRepo, authManager, hashManager, emailService)
	userService := user.NewUserService(userPostgresRepo, userCacheRepo, authPostgresRepo, objStorageRepo)
	articleService := article.NewArticleService(articlePostgresRepo, categoryPostgresRepo)
	commentService := comment.NewCommentService(commentPostgresRepo)
	categoryService := category.NewCategoryService(categoryPostgresRepo)
	likeService := like.NewLikeService(likePostgresRepo, userPostgresRepo)
	bookmarkService := bookmark.NewBookmarkService(bookmarkPostgresRepo, userPostgresRepo)

	serverDeps := server.NewServerDependencies(
		authManager,
		authService,
		userService,
		articleService,
		commentService,
		categoryService,
		likeService,
		bookmarkService,
		logger,
	)

	err = server.InitServer(config.Server.Port, serverDeps)
	checkError(logger, "Server", err)
}
