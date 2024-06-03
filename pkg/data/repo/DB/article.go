package db

import (
	"blog/pkg/data/database"
	"blog/pkg/data/models"
	"errors"

	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ArticleRepo struct {
	DB  *gorm.DB
	RDB *redis.Client
}

var (
	ErrArticleNotFound = errors.New("article not found")
)

func NewArticleRepo() *ArticleRepo {
	return &ArticleRepo{
		DB:  database.DB,
		RDB: database.Rdb,
	}
}
func (ar *ArticleRepo) GetAll() ([]map[string]string, error) {
	var articles []map[string]string
	keys, err := ar.RDB.Keys(context.Background(), "article:*").Result()
	if err != nil {
		return articles, err
	}
	for _, key := range keys {
		articleMap, err := ar.RDB.HGetAll(context.Background(), key).Result()
		if err != nil {
			return []map[string]string{}, err
		}
		articles = append(articles, articleMap)
	}
	return articles, nil
}

func (ar *ArticleRepo) GetById(id string) (map[string]string, error) {
	exists := ar.RDB.Exists(context.Background(), fmt.Sprintf("article:%s", id))
	if exists.Val() == 0 {
		return map[string]string{}, ErrArticleNotFound
	}
	redisMapRes := ar.RDB.HGetAll(context.Background(), fmt.Sprintf("article:%s", id))
	if redisMapRes.Err() != nil {
		return map[string]string{}, redisMapRes.Err()
	}
	return redisMapRes.Val(), nil
}
func (ar *ArticleRepo) Create(sAuthorId,title, content string) (models.Article, error) {
	iAuthorId,_:=strconv.Atoi(sAuthorId)
	var a models.Article
	a.Title = title
	a.Content = content
	a.AuthorId = uint(iAuthorId)
	err := ar.DB.Create(&a).Error
	if err != nil {
		return a, err
	}
	return ar.CreateChacheById(a)
}
func (ar *ArticleRepo) UpdateById(id, title, content string) (models.Article, error) {
	var a models.Article
	err := ar.DB.First(&a, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return a, ErrArticleNotFound
		}
		return a, err
	}
	a.Title = title
	a.Content = content
	err = ar.DB.Save(&a).Error
	if err != nil {
		return models.Article{}, err
	}
	err = ar.deleteChacheById(id)
	if err != nil {
		return models.Article{}, err
	}
	a, err = ar.CreateChacheById(a)
	if err != nil {
		return models.Article{}, err
	}
	return a, err
}
func (ar *ArticleRepo) DeleteById(id string) error {
	var a models.Article
	err := ar.DB.Delete(&a, id).Error
	if err != nil {
		return err
	}
	err = ar.deleteChacheById(id)
	return err
}

func (ar *ArticleRepo) CreateChacheById(a models.Article) (models.Article, error) {
	redisRes := database.Rdb.HMSet(context.Background(), fmt.Sprintf("article:%d", a.Id), map[string]interface{}{
		"title":     a.Title,
		"content":   a.Content,
		"createdAt": a.CreatedAt,
		"updatedAt": a.UpdatedAt,
		"authorId":  a.AuthorId,
	})
	return a, redisRes.Err()
}
func (ar *ArticleRepo) deleteChacheById(id string) error {
	redisRes := database.Rdb.Del(context.Background(), fmt.Sprintf("article:%s", id))
	return redisRes.Err()
}
