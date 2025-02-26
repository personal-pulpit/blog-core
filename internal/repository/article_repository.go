package repository

import (
	"blog/internal/model"
	"context"
)

type ArticlePostgresRepository interface {
	GetAll(ctx context.Context) ([]model.Article, error)
	SearchArticles(ctx context.Context,articleFilter *model.ArticleFilter) ([]model.Article, error)
	GetArticleByID(ctx context.Context,ID uint) (*model.Article, error)
	Create(ctx context.Context,articleModel *model.Article) (*model.Article, error)
	UpdateByID(ctx context.Context,ID uint, title, content string) (*model.Article, error)
	DeleteByID(ctx context.Context,ID uint) error
}
type ArticleRedisRepository interface {
	GetCaches() ([]map[string]string, error)
	GetCacheByID(ID model.ID) (map[string]string, error)
	CreateCache(ID  uint,title,content,createdAt,updatedAt string,athurID uint) error
	DeleteCacheByID(ID string) error
}