package repository

import (
	"blog/internal/model"
)

type ArticlePostgresRepository interface {
	Create(authorID model.ID, title, content string) (model.Article, error)
	UpdateByID(ID, title, content string) (model.Article, error)
	DeleteByID(ID string) error
}
type ArticleRedisRepository interface {
	GetCaches() ([]map[string]string, error)
	GetCacheByID(ID model.ID) (map[string]string, error)
	CreateCache(ID  uint,title,content,createdAt,updatedAt string,athurID uint) error
	DeleteCacheByID(ID string) error
}