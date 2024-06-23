package redis_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type articleRedisRepo struct {
	redisClient *redis.Client
}

func NewArticleRedisRepository(redisCLI *redis.Client) repository.ArticleRedisRepository {
	return &articleRedisRepo{
		redisClient: redisCLI,
	}
}
func (a *articleRedisRepo) GetCaches() ([]map[string]string, error) {
	var articles []map[string]string
	keys, err := a.redisClient.Keys(context.Background(), "article:*").Result()
	if err != nil {
		return articles, err
	}
	for _, key := range keys {
		articleMap, err := a.redisClient.HGetAll(context.Background(), key).Result()
		if err != nil {
			return []map[string]string{}, err
		}
		articles = append(articles, articleMap)
	}
	return articles, nil
}
func (a *articleRedisRepo) GetCacheByID(ID model.ID) (map[string]string, error) {
	exists := a.redisClient.Exists(context.Background(), fmt.Sprintf("article:%s", ID))
	if exists.Val() == 0 {
		return map[string]string{}, ErrArticleNotFound
	}
	redisMapRes := a.redisClient.HGetAll(context.Background(), fmt.Sprintf("article:%s", ID))
	if redisMapRes.Err() != nil {
		return redisMapRes.Val(), redisMapRes.Err()
	}
	return redisMapRes.Val(), nil
}
func (a *articleRedisRepo) CreateCache(ID uint, title, content, createdAt, updatedAt string, athurID uint) error {
	redisRes := a.redisClient.HMSet(context.Background(), fmt.Sprintf("article:%d", ID), map[string]interface{}{
		"title":     title,
		"content":   content,
		"createdAt": createdAt,
		"updatedAt": updatedAt,
		"authorId":  athurID,
	})
	return redisRes.Err()
}
func (a *articleRedisRepo) DeleteCacheByID(ID string) error {
	redisRes := a.redisClient.Del(context.Background(), fmt.Sprintf("article:%s", ID))
	return redisRes.Err()
}
