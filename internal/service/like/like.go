package like

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"errors"
)

type LikeService interface {
	CreateLike(ctx context.Context, userID, articleID uint) (*model.Like, error)
	GetUserLikes(ctx context.Context, userID uint) ([]model.Like, error)
	DeleteLike(ctx context.Context, articleID, userID uint) error
}

var ErrPermissionDenied = errors.New("permission denied for you")

type likeServiceImpl struct {
	likeRepo repository.LikeRepository
	userRepo repository.UserRepository
}

func NewLikeService(likeRepo repository.LikeRepository, userRepo repository.UserRepository) LikeService {
	return &likeServiceImpl{
		likeRepo: likeRepo,
		userRepo: userRepo,
	}
}

func (s *likeServiceImpl) CreateLike(ctx context.Context, userID, articleID uint) (*model.Like, error) {
	like := model.NewLike(userID, articleID)

	return s.likeRepo.Create(ctx, like)
}

func (s *likeServiceImpl) GetUserLikes(ctx context.Context, userID uint) ([]model.Like, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user.Likes, nil
}

func (s *likeServiceImpl) DeleteLike(ctx context.Context, articleID, userID uint) error {
	return s.likeRepo.DeleteByArticleIDAndUserID(ctx, articleID, userID)
}
