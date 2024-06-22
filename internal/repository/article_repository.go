package repository

import (
	"blog/database/mysql_repo"
	"blog/internal/model"
)

type ArticleRepo interface {
	GetAll() ([]map[string]string, error)
	GetById(id string) (map[string]string, error)
	Create(sAuthorId, title, content string) (model.Article, error)
	UpdateById(id, title, content string) (model.Article, error)
	DeleteById(id string) error
}

func NewArticleRepo() ArticleRepo {
	return mysql_repository.NewArticleRepo()
}
