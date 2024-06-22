package repository

import (
	"blog/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll() ([]map[string]string, error)
	GetByID(ID string) (map[string]string, error)
	Verify(username, password string) (model.User, error)
	Create(firstname, lastname, biography, username, password, email, phonenumber string) (model.User, *gorm.DB, error)
	UpdateByID(ID, firstname, lastname, biography, username string) (model.User, error)
	DeleteByID(ID string) error
	GetUsernameById(ID string) (string, error)
	DeleteChacheById(ID string) error
}
