package database

import (
	"blog/config"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisInstance *redis.Client
	redisMutex    = &sync.Mutex{}
)

func GetRedisDB(cfg *config.Redis) (*redis.Client,error){
	redisMutex.Lock()
	defer redisMutex.Unlock()

	if redisInstance == nil {
		url := fmt.Sprintf("redis://%s:%s@%s:%d/%d?protocol=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		cfg.Protocol,
		)
		opts, err := redis.ParseURL(url)
		if err != nil {
			return nil,err
		}
		redisInstance = redis.NewClient(opts)
	}

	return redisInstance,nil
}
func CloseRedis(){
	err := redisInstance.Close()
	if err != nil {
		panic(err)
	}
}
