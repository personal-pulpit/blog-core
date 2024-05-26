package database
import "github.com/redis/go-redis/v9"
var Rdb *redis.Client
func ConnectRedis(){
	rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
        Protocol: 3, // specify 2 for RESP 2 or 3 for RESP 3
    })
	Rdb = rdb
}
