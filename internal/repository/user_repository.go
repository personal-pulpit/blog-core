package repository

import (
	"blog/internal/model"

	"gorm.io/gorm"
)

type UserPostgresRepository interface {
	Create(*model.User) (*model.User, *gorm.DB, error)
	UpdateByID(ID, firstName, lastName, biography string) (*model.User, error)
	DeleteByID(ID model.ID) error
	GetUserByEmail(email string) (*model.User, error)
	GetUserByID(ID model.ID) (*model.User, error)
}
type UserRedisRepository interface {
	CreateCache(ID model.ID, firstName, lastName, biography, email string, role model.Role, createdAt, updatedAt string) error
	GetCacheByID(ID string) (map[string]string, error)
	DeleteCacheByID(ID string) error
}
