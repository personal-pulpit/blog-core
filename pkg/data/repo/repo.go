package repo

import (
	db "blog/pkg/data/repo/DB"
	"blog/internal/model"
	"gorm.io/gorm"
)

type UserDB interface {
	GetAll() ([]map[string]string, error)
	GetById(id string) (map[string]string, error)
	Verify(username, password string) (model.User, error)
	Create(firstname, lastname, biography, username, password, email, phonenumber string) (model.User, *gorm.DB, error)
	UpdateById(id, firstname, lastname, biography, username string) (model.User, error)
	DeleteById(id string) error
	GetUsernameById(id string) (string, error)
	DeleteChacheById(id string) error
}
type ArticleDB interface {
	GetAll() ([]map[string]string, error)
	GetById(id string) (map[string]string, error)
	Create(sAuthorId, title, content string) (model.Article, error)
	UpdateById(id, title, content string) (model.Article, error)
	DeleteById(id string) error
}

func NewUserDB() UserDB {
	return db.NewUserRepo()
}
func NewArticleDB() ArticleDB {
	return db.NewArticleRepo()
}
