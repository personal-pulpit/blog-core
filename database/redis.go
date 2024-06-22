package database

import (
	"blog/config"
	"blog/pkg/logging"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func ConnectRedis() {
	url := fmt.Sprintf("redis://%s:%s@%s:%s/%s?protocol=%s",
		config.Cfg.Redis.Username,
		config.Cfg.Redis.Password,
		config.Cfg.Redis.Host,
		config.Cfg.Redis.Port,
		config.Cfg.Redis.DBname,
		config.Cfg.Redis.Protocol,
	)
	opts, err := redis.ParseURL(url)
	if err != nil {
		logging.MyLogger.Fatal(logging.General, logging.Startup, err.Error(), nil)
	}
	RDB = redis.NewClient(opts)
}
func CloseRedis() {
	err := RDB.Close()
	if err != nil {
		logging.MyLogger.Fatal(logging.General, logging.Down, err.Error(), nil)
	}
}
