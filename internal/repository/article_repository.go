package repository

import (
	"blog/internal/model"
	"context"
)

type ArticleRepository interface {
	GetAll(ctx context.Context) ([]model.Article, error)
	SearchArticles(ctx context.Context,articleFilter *model.ArticleFilter) ([]model.Article, error)
	GetArticleByID(ctx context.Context,ID uint) (*model.Article, error)
	Create(ctx context.Context,articleModel *model.Article) (*model.Article, error)
	UpdateByID(ctx context.Context,ID uint, title, content string) (*model.Article, error)
	DeleteByID(ctx context.Context,ID uint) error
}