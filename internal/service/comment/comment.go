package comment

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
)

type CommentService interface {
	AddComment(ctx context.Context, content string, userID, articleID uint) (*model.Comment, error)
	UpdateComment(ctx context.Context, userID,commentID uint, content string) (*model.Comment, error)
	DeleteComment(ctx context.Context, userID,commentID uint) error
}

type commentServiceImpl struct {
	commentRepo repository.CommentRepository
}

func NewCommentService(commentRepository repository.CommentRepository) CommentService {
	return &commentServiceImpl{
		commentRepo: commentRepository,
	}
}

func (s *commentServiceImpl) AddComment(ctx context.Context, content string, userID, articleID uint) (*model.Comment, error) {
	comment := model.NewComment(userID, articleID, content)

	comment, err := s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *commentServiceImpl) UpdateComment(ctx context.Context, userID,commentID uint, content string) (*model.Comment, error) {
	err := s.checkCommentCreator(ctx, userID, commentID)
	if err != nil {
		return nil,err
	}

	comment, err := s.commentRepo.UpdateByID(ctx, commentID, content)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentServiceImpl) DeleteComment(ctx context.Context, userID,commentID uint) error {
	err := s.checkCommentCreator(ctx, userID, commentID)
	if err != nil {
		return err
	}

	err = s.commentRepo.DeleteByID(ctx, commentID)
	if err != nil {
		return err
	}

	return nil
}

func (s *commentServiceImpl) checkCommentCreator(ctx context.Context, userID, commentID uint) error {
	comment, err := s.commentRepo.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if comment.UserID != userID {
		return ErrPermissionDenied
	}
	return nil
}