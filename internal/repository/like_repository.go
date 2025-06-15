package repository

import (
	"blog/internal/model"
	"context"
)

type LikeRepository interface {
	Create(ctx context.Context, like *model.Like) (*model.Like, error)
	GetByID(ctx context.Context, ID uint) (*model.Like, error)
	GetByUserID(ctx context.Context, userID uint) ([]model.Like, error)
	DeleteByArticleIDAndUserID(ctx context.Context, articleID, userID uint) error
}

