package main

import (
	"blog/api/server"
	"blog/config"
	postgres "blog/database/postgres"
	redis "blog/database/redis"
	"blog/pkg/logging"
)

func main() {
	config := config.GetConfigInstance()
	logging.InitZapLogger(config.Logger)
	PostgresCLI := postgres.GetPostgresqlDB(config.Postgres)
	redisCLI := redis.GetRedisDB(config.Redis)
	defer redis.CloseRedis()
	server.InitServer(config, PostgresCLI, redisCLI)
}
