package database

import (
	"blog/config"
	"blog/pkg/logging"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisInstance *redis.Client
	redisMutex    = &sync.Mutex{}
)

func GetRedisDB() *redis.Client{
	redisMutex.Lock()
	defer redisMutex.Unlock()
	if redisInstance == nil {
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
		redisInstance = redis.NewClient(opts)
	}
	return redisInstance
}
func CloseRedis() {
	err := redisInstance.Close()
	if err != nil {
		logging.MyLogger.Fatal(logging.General, logging.Down, err.Error(), nil)
	}
}
