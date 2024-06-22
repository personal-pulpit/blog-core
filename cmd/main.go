package main

import (
	"blog/api/server"
	"blog/config"
	"blog/database"
	"blog/pkg/logging"
)

func main() {
	config.InitConfig()
	logging.InitZapLogger()
	logging.MyLogger.Info(logging.General, logging.Initialized, "configs initialized!", nil)
	logging.MyLogger.Info(logging.General, logging.Initialized, "logger initialized!", nil)
	database.ConnectDB()
	logging.MyLogger.Info(logging.General, logging.Startup, "database connected!", nil)
	database.ConnectRedis()
	logging.MyLogger.Info(logging.General, logging.Startup, "redis connected!", nil)
	defer database.CloseRedis()
	server.InitServer()
}
