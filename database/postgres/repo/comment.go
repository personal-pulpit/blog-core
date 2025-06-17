package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type CommentPostgresRepo struct {
	postgresCLI *gorm.DB
}

func NewCommentPostgresRepository(postgresCLI *gorm.DB) *CommentPostgresRepo {
	return &CommentPostgresRepo{
		postgresCLI: postgresCLI,
	}
}

func (c *CommentPostgresRepo) GetAll(ctx context.Context) ([]model.Comment, error) {
	var comments = []model.Comment{}
	err := c.postgresCLI.WithContext(ctx).Find(&comments).Error
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (c *CommentPostgresRepo) GetByID(ctx context.Context, ID uint) (*model.Comment, error) {
	comment := &model.Comment{}
	err := c.postgresCLI.WithContext(ctx).First(comment, ID).Error
	if err != nil {
		return nil, fmt.Errorf("get all comments: %w: %v", repository.ErrDatabase, err)
	}

	return comment, nil
}

func (c *CommentPostgresRepo) Create(ctx context.Context, comment *model.Comment) (*model.Comment, error) {
	err := c.postgresCLI.WithContext(ctx).Create(comment).Error
	if err != nil {
		return nil, fmt.Errorf("create comment: %v \n%w: %v", comment, repository.ErrDatabase, err)
	}

	return comment, nil
}

func (c *CommentPostgresRepo) UpdateByID(ctx context.Context, ID uint, content string) (*model.Comment, error) {
	comment := &model.Comment{}
	comment, err := c.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	comment.Content = content

	err = c.postgresCLI.WithContext(ctx).Save(comment).Error
	if err != nil {
		return nil, fmt.Errorf("update article by id: article ID:%d ,content:%s\n%w: %v", ID, content, repository.ErrDatabase, err)
	}

	return comment, nil
}

func (c *CommentPostgresRepo) DeleteByID(ctx context.Context, ID uint) error {
	comment := &model.Comment{}

	tx := c.postgresCLI.WithContext(ctx).Delete(comment, ID)
	if tx.Error != nil {
		return fmt.Errorf("delete comment by id: comment ID:%d \n%w: %v", ID, repository.ErrDatabase, tx.Error)
	}

	if tx.RowsAffected == 0 {
		return repository.ErrCommentNotFound
	}

	return nil
}
