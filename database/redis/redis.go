package database

import (
	"blog/config"
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

var (
	redisInstance *redis.Client
	redisMutex    = &sync.Mutex{}

	ErrConnectionFailed = errors.New("connection failed")
)

func GetRedisDB(cfg *config.Redis) (*redis.Client, error) {
	redisMutex.Lock()
	defer redisMutex.Unlock()

	if redisInstance == nil {
		url := fmt.Sprintf("redis://%s:%s@%s:%d/%d",
			cfg.Username,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.DB,
		)

		opts, err := redis.ParseURL(url)
		if err != nil {
			return nil, err
		}

		redisInstance = redis.NewClient(opts)

		err = redisInstance.Ping(context.TODO()).Err()
		if err != nil {
			return nil, errors.Join(err, ErrConnectionFailed)
		}

	}

	return redisInstance, nil
}
func CloseRedis() {
	err := redisInstance.Close()
	if err != nil {
		panic(err)
	}
}
