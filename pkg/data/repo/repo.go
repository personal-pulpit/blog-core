package repo

import (
	"blog/pkg/data/models"
	db "blog/pkg/data/repo/DB"

	"gorm.io/gorm"
)

type UserDB interface {
	GetAll() ([]map[string]string, error)
	GetById(id string) (map[string]string, error)
	Verify(username, password string) (models.User, error)
	Create(firstname, lastname, biography, username, password, email, phonenumber string) (models.User,*gorm.DB,error)
	UpdateById(id, firstname, lastname, biography, username string) (models.User, error)
	DeleteById(id string) error
	GetUsernameById(id string)(string,error)
	DeleteChacheById(id string) error
}
type ArticleDB interface {
	GetAll() ([]map[string]string, error)
	GetById(id string) (map[string]string, error)
	Create(sAuthorId,title, content string) (models.Article, error)	
	UpdateById(id, title, content string) (models.Article, error) 
	DeleteById(id string) error
}
func NewUserDB()UserDB{
	return db.NewUserRepo()
}
func NewArticleDB()ArticleDB{
	return db.NewArticleRepo()
}
