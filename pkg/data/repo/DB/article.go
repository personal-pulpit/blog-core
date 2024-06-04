package db

import (
	"blog/pkg/data/database"
	"blog/pkg/data/models"
	"blog/pkg/logging"
	"errors"

	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ArticleRepo struct {
	DB     *gorm.DB
	RDB    *redis.Client
	Logger logging.ZapLogger
}

var (
	ErrArticleNotFound = errors.New("article not found")
)

func NewArticleRepo() *ArticleRepo {
	return &ArticleRepo{
		DB:     database.DB,
		RDB:    database.Rdb,
		Logger: logging.MyLogger,
	}
}
func (ar *ArticleRepo) GetAll() ([]map[string]string, error) {
	var articles []map[string]string
	keys, err := ar.RDB.Keys(context.Background(), "article:*").Result()
	if err != nil {
		ar.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
		return articles, err
	}
	for _, key := range keys {
		articleMap, err := ar.RDB.HGetAll(context.Background(), key).Result()
		if err != nil {
			ar.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
			return []map[string]string{}, err
		}
		articles = append(articles, articleMap)
	}
	ar.Logger.Info(logging.Redis, logging.Get, "", nil)
	return articles, nil
}

func (ar *ArticleRepo) GetById(id string) (map[string]string, error) {
	exists := ar.RDB.Exists(context.Background(), fmt.Sprintf("article:%s", id))
	if exists.Val() == 0 {
		ar.Logger.Error(logging.Redis, logging.Get, ErrArticleNotFound.Error(), nil)
		return map[string]string{}, ErrArticleNotFound
	}
	redisMapRes := ar.RDB.HGetAll(context.Background(), fmt.Sprintf("article:%s", id))
	if redisMapRes.Err() != nil {
		ar.Logger.Error(logging.Redis, logging.Get, redisMapRes.Err().Error(), nil)
		return redisMapRes.Val(), redisMapRes.Err()
	}
	ar.Logger.Info(logging.Redis, logging.Get, "", nil)
	return redisMapRes.Val(), nil
}
func (ar *ArticleRepo) Create(sAuthorId, title, content string) (models.Article, error) {
	iAuthorId, _ := strconv.Atoi(sAuthorId)
	var a models.Article
	a.Title = title
	a.Content = content
	a.AuthorId = uint(iAuthorId)
	tx := NewTx(ar.DB)
	err := tx.Create(&a).Error
	if err != nil {
		ar.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
		return a, err
	}
	ar.Logger.Info(logging.Mysql, logging.Insert, "", nil)
	err = ar.CreateChacheById(a)
	if err != nil{
		tx.Rollback()
		ar.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return a,err
	}
	tx.Commit()
	ar.Logger.Error(logging.Mysql, logging.Insert, "", nil)
	return a,nil
}
func (ar *ArticleRepo) UpdateById(id, title, content string) (models.Article, error) {
	var a models.Article
	err := ar.DB.First(&a, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ar.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
			return a, ErrArticleNotFound
		}
		ar.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return a, err
	}
	a.Title = title
	a.Content = content
	tx := NewTx(ar.DB)
	err = tx.Save(&a).Error
	if err != nil {
		ar.Logger.Error(logging.Mysql, logging.Update, err.Error(), nil)
		return a, err
	}
	err = ar.CreateChacheById(a)
	if err != nil {
		tx.Rollback()
		ar.Logger.Error(logging.Redis, logging.Set, err.Error(), nil)
		ar.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return a, err
	}
	tx.Commit()
	ar.Logger.Info(logging.Mysql, logging.Insert, "", nil)
	return a, err
}
func (ar *ArticleRepo) DeleteById(id string) error {
	var a models.Article
	tx := NewTx(ar.DB)
	err := tx.Delete(&a, id).Error
	if err != nil {
		ar.Logger.Error(logging.Mysql, logging.Delete, err.Error(), nil)
		return err
	}
	err = ar.deleteChacheById(id)
	if err != nil{
		tx.Rollback()
		ar.Logger.Error(logging.Redis, logging.Delete, err.Error(), nil)
		ar.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return err
	}
	ar.Logger.Info(logging.Mysql, logging.Delete, "", nil)
	return nil
}

func (ar *ArticleRepo) CreateChacheById(a models.Article)  error {
	redisRes := database.Rdb.HMSet(context.Background(), fmt.Sprintf("article:%d", a.Id), map[string]interface{}{
		"title":     a.Title,
		"content":   a.Content,
		"createdAt": a.CreatedAt,
		"updatedAt": a.UpdatedAt,
		"authorId":  a.AuthorId,
	})
	return redisRes.Err()
}
func (ar *ArticleRepo) deleteChacheById(id string) error {
	redisRes := database.Rdb.Del(context.Background(), fmt.Sprintf("article:%s", id))
	return redisRes.Err()
}
