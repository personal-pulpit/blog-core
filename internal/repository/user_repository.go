package repository

import (
	"blog/internal/model"

	"gorm.io/gorm"
)

type UserMysqlRepository interface {
	Create(*model.User) (*model.User, *gorm.DB, error)
	UpdateByID(ID, firstName, lastName, biography string) (*model.User, error)
	DeleteByID(ID model.ID) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(ID model.ID) (*model.User, error)
}
type UserRedisRepository interface {
	CreateCache(ID uint, firstName, lastName, biography, username, email, phonenumber string, role int, createdAt, updatedAt string) error
	GetCaches() ([]map[string]string, error)
	GetCacheByID(ID string) (map[string]string, error)
	DeleteCacheByID(ID string) error
}
