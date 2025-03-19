package repository

import (
	"blog/internal/model"
	"context"
)

type CategoryRepository interface {
	CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error)
	GetAllCategories(ctx context.Context) ([]model.Category, error)
	GetCategoryByID(ctx context.Context, categoryID uint) (*model.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*model.Category, error)

	DeleteCategoryByID(ctx context.Context, categoryID uint) error
}
