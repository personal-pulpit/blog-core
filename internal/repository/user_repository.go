package repository

import (
	"blog/database/mysql_repo"
	"blog/internal/model"

	"gorm.io/gorm"

)

type UserRepo interface {
	GetAll() ([]map[string]string, error)
	GetById(id string) (map[string]string, error)
	Verify(username, password string) (model.User, error)
	Create(firstname, lastname, biography, username, password, email, phonenumber string) (model.User, *gorm.DB, error)
	UpdateById(id, firstname, lastname, biography, username string) (model.User, error)
	DeleteById(id string) error
	GetUsernameById(id string) (string, error)
	DeleteChacheById(id string) error
}

func NewUserRepo() UserRepo {
	return mysql_repository.NewUserRepo()
}
