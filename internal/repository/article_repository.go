package repository

import (
	"blog/internal/model"
)

type ArticleRepository interface {
	GetAll() ([]map[string]string, error)
	GetByID(ID string) (map[string]string, error)
	Create(sAuthorId, title, content string) (model.Article, error)
	UpdateByID(ID, title, content string) (model.Article, error)
	DeleteByID(ID string) error
}
