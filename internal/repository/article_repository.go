package repository

import (
	"blog/internal/model"
)

type ArticlePostgresRepository interface {
	GetAll() ([]*model.Article, error)
	GetArticle(filters map[string]interface{}) (*model.Article, error)
	GetArticleByTitle(title string) (*model.Article, error)
	GetArticleById(id model.ID) (*model.Article, error)
	Create(articleModel *model.Article) (*model.Article, error)
	UpdateByID(ID model.ID, title, content string) (*model.Article, error)
	DeleteByID(ID model.ID) error
}
type ArticleRedisRepository interface {
	GetCaches() ([]map[string]string, error)
	GetCacheByID(ID model.ID) (map[string]string, error)
	CreateCache(ID  uint,title,content,createdAt,updatedAt string,athurID uint) error
	DeleteCacheByID(ID string) error
}