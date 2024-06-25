package redis_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type userRedisRepo struct {
	redisClient *redis.Client
}

func NewUserRedisRepository(redisCLI *redis.Client) repository.UserRedisRepository {
	return &userRedisRepo{
		redisClient: redisCLI,
	}
}

func (u *userRedisRepo) CreateCache(ID model.ID, FirstName, lastname, biography, email string, role model.Role, createdAt, updatedAt string) error {
	redisRes := u.redisClient.HMSet(context.Background(), fmt.Sprintf("user:%s", ID), map[string]interface{}{
		"firstName":   FirstName,
		"lastName":    lastname,
		"biography":   biography,
		"email":       email,
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
