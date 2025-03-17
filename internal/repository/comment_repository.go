package repository

import (
	"blog/internal/model"
	"context"
)

type CommentRepository interface {
	GetAll(ctx context.Context) ([]model.Comment, error)
	GetByID(ctx context.Context,ID uint) (*model.Comment, error)
	Create(ctx context.Context,commentModel *model.Comment) (*model.Comment, error)
	UpdateByID(ctx context.Context,ID uint, content string) (*model.Comment, error)
	DeleteByID(ctx context.Context,ID uint) error
}
