package main

import (
	"blog/api/server"
	"blog/config"
	"blog/pkg/data/database"
	"blog/pkg/logger"
)

func main() {
	logger.InitZapLogger()
	logger.MyLogger.Info("logger initialized!",map[string]interface{}{
		"status":true,
	})
	config.InitConfig()
	logger.MyLogger.Info("configs initialized!",map[string]interface{}{
		"status":true,
	})
	database.ConnectDB()
	logger.MyLogger.Info("database connected!",map[string]interface{}{
		"status":true,
	})
	server.InitServer()
}