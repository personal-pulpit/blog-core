package repository

import (
	"blog/internal/model"
	"context"
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) (*model.User, *gorm.DB, error)
	UpdateByID(ctx context.Context, ID uint, firstName, lastName, biography string) (*model.User, error)
	DeleteByID(ctx context.Context, ID uint) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserByID(ctx context.Context, ID uint) (*model.User, error)
}

type UserCacheRepository interface {
	SetUserProfileImageURL(ctx context.Context, userID uint, url string, exp time.Duration) error
	GetUserProfileImageURL(ctx context.Context, userID uint) (string, error)
	DestroyUserProfileImageURL(ctx context.Context, userID uint) error
}
