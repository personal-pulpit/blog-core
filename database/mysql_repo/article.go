package mysql_repository

import (
	"blog/database"
	"blog/internal/model"
	"blog/internal/repository"
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

func NewArticleRepo() repository.ArticleRepository {
	return &ArticleRepo{
		DB:     database.GetMysqlDB(),
		RDB:    database.GetRedisDB(),
		Logger: logging.MyLogger,
	}
}
func (a *ArticleRepo) GetAll() ([]map[string]string, error) {
	var articles []map[string]string
	keys, err := a.RDB.Keys(context.Background(), "article:*").Result()
	if err != nil {
		a.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
		return articles, err
	}
	for _, key := range keys {
		articleMap, err := a.RDB.HGetAll(context.Background(), key).Result()
		if err != nil {
			a.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
			return []map[string]string{}, err
		}
		articles = append(articles, articleMap)
	}
	a.Logger.Info(logging.Redis, logging.Get, "", nil)
	return articles, nil
}

func (a *ArticleRepo) GetByID(ID string) (map[string]string, error) {
	exists := a.RDB.Exists(context.Background(), fmt.Sprintf("article:%s", ID))
	if exists.Val() == 0 {
		a.Logger.Error(logging.Redis, logging.Get, ErrArticleNotFound.Error(), nil)
		return map[string]string{}, ErrArticleNotFound
	}
	redisMapRes := a.RDB.HGetAll(context.Background(), fmt.Sprintf("article:%s", ID))
	if redisMapRes.Err() != nil {
		a.Logger.Error(logging.Redis, logging.Get, redisMapRes.Err().Error(), nil)
		return redisMapRes.Val(), redisMapRes.Err()
	}
	a.Logger.Info(logging.Redis, logging.Get, "", nil)
	return redisMapRes.Val(), nil
}
func (a *ArticleRepo) Create(sAuthorId, title, content string) (model.Article, error) {
	iAuthorId, _ := strconv.Atoi(sAuthorId)
	var article model.Article
	article.Title = title
	article.Content = content
	article.AuthorId = uint(iAuthorId)
	tx := NewTx(a.DB)
	err := tx.Create(&article).Error
	if err != nil {
		a.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
		return article, err
	}
	a.Logger.Info(logging.Mysql, logging.Insert, "", nil)
	err = a.CreateChacheById(article)
	if err != nil {
		tx.Rollback()
		a.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return article, err
	}
	tx.Commit()
	a.Logger.Error(logging.Mysql, logging.Insert, "", nil)
	return article, nil
}
func (a *ArticleRepo) UpdateByID(ID, title, content string) (model.Article, error) {
	var article model.Article
	err := a.DB.First(&article, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			a.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
			return article, ErrArticleNotFound
		}
		a.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return article, err
	}
	article.Title = title
	article.Content = content
	tx := NewTx(a.DB)
	err = tx.Save(&article).Error
	if err != nil {
		a.Logger.Error(logging.Mysql, logging.Update, err.Error(), nil)
		return article, err
	}
	err = a.CreateChacheById(article)
	if err != nil {
		tx.Rollback()
		a.Logger.Error(logging.Redis, logging.Set, err.Error(), nil)
		a.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return article, err
	}
	tx.Commit()
	a.Logger.Info(logging.Mysql, logging.Insert, "", nil)
	return article, err
}
func (a *ArticleRepo) DeleteByID(ID string) error {
	var article model.Article
	tx := NewTx(a.DB)
	err := tx.Delete(&article, ID).Error
	if err != nil {
		a.Logger.Error(logging.Mysql, logging.Delete, err.Error(), nil)
		return err
	}
	err = a.deleteChacheById(ID)
	if err != nil {
		tx.Rollback()
		a.Logger.Error(logging.Redis, logging.Delete, err.Error(), nil)
		a.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return err
	}
	a.Logger.Info(logging.Mysql, logging.Delete, "", nil)
	return nil
}

func (a *ArticleRepo) CreateChacheById(article model.Article) error {
	redisRes := a.RDB.HMSet(context.Background(), fmt.Sprintf("article:%d", article.ID), map[string]interface{}{
		"title":     article.Title,
		"content":   article.Content,
		"createdAt": article.CreatedAt,
		"updatedAt": article.UpdatedAt,
		"authorId":  article.AuthorId,
	})
	return redisRes.Err()
}
func (a *ArticleRepo) deleteChacheById(ID string) error {
	redisRes := a.RDB.Del(context.Background(), fmt.Sprintf("article:%s", ID))
	return redisRes.Err()
}
