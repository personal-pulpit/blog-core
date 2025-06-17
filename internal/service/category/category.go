package category

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, name string) (*model.Category, error)
	GetAllCategories(ctx context.Context) ([]model.Category, error)
	GetCategoryByID(ctx context.Context, categoryID uint) (*model.Category, error)
	GetCategoryByTitle(ctx context.Context, title string) (*model.Category, error)
	DeleteCategoryByID(ctx context.Context, categoryID uint) error
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *categoryService) CreateCategory(ctx context.Context, name string) (*model.Category, error) {
	category := model.NewCategory(name)

	return s.categoryRepo.CreateCategory(ctx, category)
}

func (s *categoryService) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	return s.categoryRepo.GetAllCategories(ctx)
}

func (s *categoryService) GetCategoryByID(ctx context.Context, categoryID uint) (*model.Category, error) {
	return s.categoryRepo.GetCategoryByID(ctx, categoryID)
}

func (s *categoryService) GetCategoryByTitle(ctx context.Context, title string) (*model.Category, error) {

	return s.categoryRepo.GetCategoryByName(ctx, title)
}

func (s *categoryService) DeleteCategoryByID(ctx context.Context, categoryID uint) error {
	return s.categoryRepo.DeleteCategoryByID(ctx, categoryID)
}
