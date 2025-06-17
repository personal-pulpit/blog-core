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

	likeID    uint
	articleID uint
	userID    uint
}

func (s *LikeTestSuite) SetupSuite() {
	s.likeRepo = postgres_repository.NewLikePostgresRepository(db)
	userRepo := postgres_repository.NewUserPostgresRepository(db)
	articleRepo := postgres_repository.NewArticlePostgresRepo(db)

	userModel := model.NewUser("user1 firstName", "user1 lastName", "sdudjkalakakd@fake.com", "user1 biography")
	user, tx, err := userRepo.Create(context.TODO(), userModel)
	s.NoError(err)
	s.NotNil(user)
	tx.Commit()

	s.userID = user.ID

	article, err := articleRepo.Create(context.TODO(), model.NewArticle("article1", "content of article1", user.ID, nil))
	s.NoError(err)
	s.NotNil(article)

	s.articleID = article.ID

	s.service = like.NewLikeService(s.likeRepo, userRepo)
}

func (s *LikeTestSuite) TestA_CreateLike() {
	ctx := context.TODO()

	like, err := s.service.CreateLike(ctx, s.userID, s.articleID)

	s.NoError(err)
	s.NotNil(like)

	s.likeID = like.ID
}

func (s *LikeTestSuite) TestB_GetUserLikes() {
	ctx := context.TODO()

	testCases := []struct {
		userID uint
		Valid  bool
	}{
		{
			userID: 0,
			Valid:  false,
		},
		{
			userID: s.userID,
			Valid:  true,
		},
	}

	for _, tc := range testCases {
		likes, err := s.service.GetUserLikes(ctx, tc.userID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(likes)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *LikeTestSuite) TestC_DeleteLIke() {
	ctx := context.TODO()

	testCases := []struct {
		userID    uint
		articleID uint
		Valid     bool
	}{
		{
			userID:    s.userID,
			articleID: 0,
			Valid:     false,
		},
		{
			userID:    0,
			articleID: s.likeID,
			Valid:     false,
		},
		{
			userID:    s.userID,
			articleID: s.articleID,
			Valid:     true,
		},
	}

	for _, tc := range testCases {
		err := s.service.DeleteLike(ctx, tc.articleID, tc.userID)
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
