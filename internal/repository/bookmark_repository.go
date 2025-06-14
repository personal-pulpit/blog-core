package repository

import (
	"blog/internal/model"
	"context"
)

type BookmarkRepository interface {
	Create(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)
	GetByUserID(ctx context.Context, userID uint) ([]*model.Bookmark, error)
	GetByID(ctx context.Context, bookmarkID uint) (*model.Bookmark, error)
	DeleteByID(ctx context.Context, bookmarkID uint) error
}