package service

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/internal/service/like"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LikeTestSuite struct {
	suite.Suite
	likeRepo repository.LikeRepository

	service like.LikeService

	likeID uint
	articleID uint
	userID    uint
}

func (s *LikeTestSuite) SetupSuite() {
	s.likeRepo = postgres_repository.NewLikePostgresRepository(db)
	userRepo := postgres_repository.NewUserPostgresRepository(db)
	articleRepo := postgres_repository.NewArticlePostgresRepo(db)

	userModel := model.NewUser("user1 firstName", "user1 lastName", "afakeonce@fake.come", "user1 biography")
	user, tx, err := userRepo.Create(context.TODO(), userModel)
	s.NoError(err)
	s.NotNil(user)
	tx.Commit()

	s.userID = user.ID

	article, err := articleRepo.Create(context.TODO(), model.NewArticle("article1", "content of article1", user.ID, nil))
	s.NoError(err)
	s.NotNil(article)

	s.articleID = article.ID

	s.service = like.NewLikeService(s.likeRepo)
}

func (s *LikeTestSuite) TestA_CreateLike() {
	ctx := context.TODO()

	like, err := s.service.CreateLike(ctx, s.userID, s.articleID)

	s.NoError(err)
	s.NotNil(like)

	s.likeID = like.ID
}

func (s *LikeTestSuite) TestC_DeleteLIke() {
	ctx := context.TODO()

	testCases := []struct {
		likeID uint
		userID uint
		Valid     bool
	}{
		{
			likeID: 0,
			userID:s.userID,
			Valid:     false,
		},
		{
			likeID: s.likeID,
			userID: 0,
			Valid:     false,
		},
		{
			likeID: s.likeID,
			userID:s.userID,
			Valid:     true,
		},
	}

	for _, tc := range testCases {
		err := s.service.DeleteLike(ctx, tc.userID,tc.likeID)
		if tc.Valid {
			s.NoError(err)

		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func TestLikeSuite(t *testing.T) {
	t.Helper()

	suite.Run(t, new(LikeTestSuite))
}
