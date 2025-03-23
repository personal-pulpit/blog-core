package redis_repo

import (
	"blog/internal/repository"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type redisRepository struct {
	redisCLI *redis.Client
}

func NewRedisRepository(redisClient *redis.Client) repository.UserCacheRepository {
	return &redisRepository{
		redisCLI: redisClient,
	}
}

func (r *redisRepository) SetUserProfileImageURL(ctx context.Context, userID uint, url string, exp time.Duration) error {
	key := fmt.Sprintf("user:%d", userID)

	err := r.redisCLI.HSet(ctx, key, "profile_image_url", url).Err()
	if err != nil {
		return fmt.Errorf("failed to set user profile image URL in redis: %w", err)
	}

	err = r.redisCLI.Expire(ctx, key, exp).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration for user profile image URL in redis: %w", err)
	}


	return nil
}

func (r *redisRepository) GetUserProfileImageURL(ctx context.Context, userID uint) (string, error) {
    key := fmt.Sprintf("user:%d", userID)

    url, err := r.redisCLI.HGet(ctx, key, "profile_image_url").Result()
    if err != nil {
        if err == redis.Nil {
            return "", repository.ErrUserImageURLNotFound
        }

        return "", fmt.Errorf("failed to get user profile image URL from redis: %w", err)
    }

    return url, nil
}

func (r *redisRepository) DestroyUserProfileImageURL(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user:%d", userID)

	err := r.redisCLI.HDel(ctx, key,"profile_image_url").Err()
	if err != nil {
		return fmt.Errorf("failed to delete user profile image URL from redis: %w", err)
	}

	return nil
}
