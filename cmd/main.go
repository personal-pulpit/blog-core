package main

import (
	"blog/api/server"
	"blog/config"
	redis "blog/database/redis"
	mysql "blog/database/mysql"
	"blog/pkg/logging"
)

func main() {
	config := config.GetConfigInstance()
	logging.InitZapLogger(config.Logger)
	mysqlCLI := mysql.GetMysqlDB(config.Mysql)
	redisCLI := redis.GetRedisDB(config.Redis)	
	defer redis.CloseRedis()
	server.InitServer(config,mysqlCLI,redisCLI)
}
