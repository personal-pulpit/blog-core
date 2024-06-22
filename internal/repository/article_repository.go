package repository

import (
	"blog/internal/model"
)

type ArticleMysqlRepository interface {
	Create(sAuthorId, title, content string) (model.Article, error)
	UpdateByID(ID, title, content string) (model.Article, error)
	DeleteByID(ID string) error
}
type ArticleRedisRepository interface {
	GetCaches() ([]map[string]string, error)
	GetCacheByID(ID string) (map[string]string, error)
	CreateCache(ID  uint,title,content,createdAt,updatedAt string,athurID uint) error
	DeleteCacheByID(ID string) error
}