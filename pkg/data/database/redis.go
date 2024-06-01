package database

import (
	"blog/config"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

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
		panic(err)
	}
	Rdb = redis.NewClient(opts)
}
func CloseRedis() {
	err := Rdb.Close()
	if err != nil{
		log.Fatalln(err)
	}
}
