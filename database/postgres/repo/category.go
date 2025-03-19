package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type categoryPostgresRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryPostgresRepository{
		db: db,
	}
}

func (c *categoryPostgresRepository) CreateCategory(ctx context.Context, category *model.Category) (*model.Category, error) {
	err := c.db.WithContext(ctx).Create(&category).Error
	if err != nil {
		return nil, fmt.Errorf("create category: %v\n%w: %v", category, repository.ErrDatabase, err)
	}

	return category, nil
}

func (c *categoryPostgresRepository) GetAllCategories(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	err := c.db.WithContext(ctx).Preload("Articles").Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("get all categories: %w: %v", repository.ErrDatabase, err)
	}

	return categories, nil
}

func (c *categoryPostgresRepository) GetCategoryByID(ctx context.Context, categoryID uint) (*model.Category, error) {
	category := new(model.Category)
	err := c.db.WithContext(ctx).Preload("Articles").First(category, categoryID).Error
	if err != nil {
		return nil, fmt.Errorf("get category by id: id:%d \n%w: %v", categoryID,repository.ErrDatabase, err)
	}

	return category, nil
}

func (c *categoryPostgresRepository) GetCategoryByName(ctx context.Context, name string) (*model.Category, error) {
	category := new(model.Category)
	err := c.db.WithContext(ctx).Where("name = ?", name).First(category).Error
	if err != nil {
		return nil, fmt.Errorf("get category by name: name:%s \n%w: %v", name,repository.ErrDatabase, err)
	}

	return category, nil
}

func (c *categoryPostgresRepository) DeleteCategoryByID(ctx context.Context, categoryID uint) error {
	err := c.db.WithContext(ctx).Delete(&model.Category{}, categoryID).Error
	if err != nil {
		return fmt.Errorf("delete category by id: id:%d \n%w: %v", categoryID,repository.ErrDatabase, err)
	}

	return nil
}




