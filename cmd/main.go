package main

import (
	"blog/api/server"
	"blog/config"
	redis "blog/database/redis"
	mysql "blog/database/mysql"
	"blog/pkg/logging"
)

func main() {
	config.InitConfig()
	logging.InitZapLogger()
	mysqlCLI := mysql.GetMysqlDB()
	redisCLI := redis.GetRedisDB()	
	defer redis.CloseRedis()
	server.InitServer(mysqlCLI,redisCLI)
}
