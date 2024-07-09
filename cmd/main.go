package main

import (
	"blog/api/server"
	"blog/config"
	postgres "blog/database/postgres"
	redis "blog/database/redis"
	"blog/pkg/logger"
)

func checkError(loggerInstance logger.Logger,err error){
	if err != nil{
		loggerInstance.Fatal(logger.General,logger.Startup,err.Error(),map[logger.ExtraKey]interface{}{})
	}
}
func main() {
	config := config.GetConfigInstance()

	logger := logger.GetZapLoggerInstance(&config.Logger)

	PostgresCLI,err := postgres.GetPostgresqlDB(&config.Postgres)
	checkError(logger,err)

	redisCLI,err := redis.GetRedisDB(&config.Redis)
	checkError(logger,err)


	defer redis.CloseRedis()

	err = server.InitServer(config, PostgresCLI, redisCLI,logger)
	checkError(logger,err)
}
