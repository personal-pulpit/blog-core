package like

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"errors"
)

type LikeService interface {
	CreateLike(ctx context.Context, userID, articleID uint) (*model.Like, error)
	DeleteLike(ctx context.Context, userID, likeID uint) error
}

var ErrPermissionDenied = errors.New("permission denied for you")

type likeServiceImpl struct {
	likeRepo repository.LikeRepository
}

func NewLikeService(likeRepo repository.LikeRepository) LikeService {
	return &likeServiceImpl{
		likeRepo: likeRepo,
	}
}

func (s *likeServiceImpl) CreateLike(ctx context.Context, userID, articleID uint) (*model.Like, error) {
	like := model.NewLike(userID, articleID)

	return s.likeRepo.Create(ctx, like)
}

func (s *likeServiceImpl) DeleteLike(ctx context.Context, userID, likeID uint) error {
	err := s.checkLikeOwner(ctx,userID,likeID)
	if err != nil {
		return err
	}

	return s.likeRepo.DeleteByID(ctx, likeID)
}

func (s *likeServiceImpl) checkLikeOwner(ctx context.Context, userID, likeID uint) error {
	like, err := s.likeRepo.GetByID(ctx, likeID)
	if err != nil {
		return err
	}

	if like.UserID != userID {
		return ErrPermissionDenied
	}
	return nil
}
