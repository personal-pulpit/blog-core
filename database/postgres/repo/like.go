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

func (l *likePostgresRepository) DeleteByID(ctx context.Context, ID uint) error {
	result := l.postgresCLI.WithContext(ctx).Delete(&model.Like{},ID)

	if result.Error != nil {
		return fmt.Errorf("delete like: %d\n%w: %v", ID, repository.ErrDatabase, result.Error)
	}

	if result.RowsAffected == 0 {
		return repository.ErrLikeNotFound
	}


	return nil
}