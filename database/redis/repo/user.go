package redis_repository

import (
	db "blog/database/redis"
	"blog/internal/repository"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type userRedisRepo struct {
	redisClient *redis.Client
}

func NewUserRedisRepository()repository.UserRedisRepository  {
	return &userRedisRepo{
		redisClient: db.GetRedisDB(),
	}
}
func (u *userRedisRepo) GetCaches() ([]map[string]string, error) {
	var users []map[string]string
	keys, err := u.redisClient.Keys(context.Background(), "user:*").Result()
	if err != nil {
		return users, err
	}
	for _, key := range keys {
		userMap, err := u.redisClient.HGetAll(context.Background(), key).Result()
		if err != nil {
			return []map[string]string{}, err
		}
		users = append(users, userMap)
	}
	return users, nil
}
func (u *userRedisRepo) CreateCache(ID uint, firstname, lastname, biography, username, email, phonenumber string, role int, createdAt, updatedAt string) error {
	redisRes := u.redisClient.HMSet(context.Background(), fmt.Sprintf("user:%d", ID), map[string]interface{}{
		"firstname":   firstname,
		"lastname":    lastname,
		"biography":   biography,
		"username":    username,
		"email":       email,
		"phonenumber": phonenumber,
		"role":        role,
		"createdAt":   createdAt,
		"updatedAt":   updatedAt,
	})
	return redisRes.Err()
}
func (u *userRedisRepo) GetCacheByID(ID string) (map[string]string, error) {
	exists := u.redisClient.Exists(context.Background(), fmt.Sprintf("user:%s", ID))
	if exists.Val() == 0 {
		return map[string]string{}, ErrUserNotFound
	}
	redisMapRes := u.redisClient.HGetAll(context.Background(), fmt.Sprintf("user:%s", ID))
	if redisMapRes.Err() != nil {
		return map[string]string{}, redisMapRes.Err()
	}
	return redisMapRes.Val(), nil
}
func (u *userRedisRepo) DeleteCacheByID(ID string) error {
	redisRes := u.redisClient.Del(context.Background(), fmt.Sprintf("user:%s", ID))
	return redisRes.Err()
}
