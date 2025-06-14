package repository

import (
	"blog/internal/model"
	"context"
)

type LikeRepository interface {
	Create(ctx context.Context, like *model.Like) (*model.Like, error)
	GetByID(ctx context.Context, ID uint) (*model.Like, error)
	DeleteByID(ctx context.Context, ID uint) error
}
