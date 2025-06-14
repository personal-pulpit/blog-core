package postgres_repository

import (
	"blog/internal/model"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type bookmarkPostgresRepository struct {
	postgresCLI *gorm.DB
}

func NewBookmarkPostgresRepository(postgresCLI *gorm.DB) *bookmarkPostgresRepository {
	return &bookmarkPostgresRepository{
		postgresCLI: postgresCLI,
	}
}

func (b *bookmarkPostgresRepository) Create(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	err := b.postgresCLI.WithContext(ctx).Create(bookmark).Error
	if err != nil {
		return nil, fmt.Errorf("create bookmark: %w", err)
	}

	return bookmark, nil
}

func (b *bookmarkPostgresRepository) GetByID(ctx context.Context, bookmarkID uint) (*model.Bookmark, error) {
	bookmark := &model.Bookmark{}

	err := b.postgresCLI.WithContext(ctx).Preload("Article").First(bookmark, bookmarkID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("bookmark not found with id: %d", bookmarkID)
		}
		return nil, fmt.Errorf("get bookmark by id: %w", err)
	}

	return bookmark, nil
}

func (b *bookmarkPostgresRepository) GetByUserID(ctx context.Context, userID uint) ([]*model.Bookmark, error) {
	var bookmarks []*model.Bookmark

	err := b.postgresCLI.WithContext(ctx).Where("user_id = ?", userID).Preload("Article").Find(&bookmarks).Error
	if err != nil {
		return nil, fmt.Errorf("get bookmarks by user id: %w", err)
	}

	return bookmarks, nil
}

func (b *bookmarkPostgresRepository) DeleteByID(ctx context.Context, bookmarkID uint) error {
	result := b.postgresCLI.WithContext(ctx).Delete(&model.Bookmark{}, bookmarkID)

	if result.Error != nil {
		return fmt.Errorf("delete bookmark: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no bookmark found with id: %d", bookmarkID)
	}

	return nil
}
