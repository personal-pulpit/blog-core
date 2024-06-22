package repository

import (
	"blog/internal/model"

	"gorm.io/gorm"
)
type UserMysqlRepository interface {
	Create(firstname, lastname, biography, username, password, email, phonenumber string) (model.User, *gorm.DB, error)
	Verify(username, password string) (model.User, error)
	UpdateByID(ID, firstname, lastname, biography, username string) (model.User, error)
	DeleteByID(ID string) error
}
type UserRedisRepository interface {
	CreateCache(ID uint, firstname, lastname, biography, username, email, phonenumber string, role int, createdAt, updatedAt string) error
	GetCaches() ([]map[string]string, error)
	GetCacheByID(ID string) (map[string]string, error)
	DeleteCacheByID(ID string) error
}
