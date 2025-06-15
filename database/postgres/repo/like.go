package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type likePostgresRepository struct {
	postgresCLI *gorm.DB
}

func NewLikePostgresRepository(db *gorm.DB) *likePostgresRepository {
	return &likePostgresRepository{
		postgresCLI: db,
	}
}

func (l *likePostgresRepository) Create(ctx context.Context, like *model.Like) (*model.Like, error) {
	err := l.postgresCLI.WithContext(ctx).Create(&like).Error
	if err != nil {
		return nil, fmt.Errorf("create like: %v\n%w: %v", like, repository.ErrDatabase, err)
	}

	return like, nil
}

func (l *likePostgresRepository) GetByID(ctx context.Context, ID uint)(*model.Like, error) {
	var like = &model.Like{}

	err := l.postgresCLI.WithContext(ctx).Preload("Article").First(like,ID).Error
	if err != nil {
		return nil, fmt.Errorf("get likes by id: %w: %v", repository.ErrDatabase, err)
	}

	return like,nil
}

func(l *likePostgresRepository) GetByUserID(ctx context.Context, userID uint) ([]model.Like, error) {
	var likes []model.Like

	err := l.postgresCLI.WithContext(ctx).Where("user_id = ?", userID).Preload("Article").Find(&likes).Error
	if err != nil {
		return nil, fmt.Errorf("get likes by user id: %w: %v", repository.ErrDatabase, err)
	}

	return likes, nil
}

func (l *likePostgresRepository) DeleteByArticleIDAndUserID(ctx context.Context, articleID,userID uint) error {
	result := l.postgresCLI.WithContext(ctx).Where("article_id = ? AND user_id = ?", articleID, userID).Delete(&model.Like{})

	if result.Error != nil {
		return fmt.Errorf("delete like: %d\n%w: %v", articleID, repository.ErrDatabase, result.Error)
	}

	if result.RowsAffected == 0 {
		return repository.ErrLikeNotFound
	}


	return nil
}