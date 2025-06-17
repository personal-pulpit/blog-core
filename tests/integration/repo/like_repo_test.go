package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type LikeTestSuite struct {
	suite.Suite
	repo        repository.LikeRepository
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
	likeID      uint
	articleID   uint
	userID      uint
}

func (s *LikeTestSuite) SetupSuite() {
	s.repo = postgres_repository.NewLikePostgresRepository(db)
	s.userRepo = postgres_repository.NewUserPostgresRepository(db)
	s.articleRepo = postgres_repository.NewArticlePostgresRepo(db)
	// categoryRepo := ostgres_repository.NewCategoryRepository(db)

	user, tx, err := s.userRepo.Create(context.TODO(), model.NewUser("user2", "user2", "<TEST-Like>", "user2"))
	s.Nil(err)
	s.NotNil(user)
	tx.Commit()

	s.userID = user.ID

	article, err := s.articleRepo.Create(context.TODO(), model.NewArticle("article1", "content", s.userID, nil))
	s.Nil(err)
	s.NotNil(article)

	s.articleID = article.ID
}

func (s *LikeTestSuite) TestA_CreateLike() {
	ctx := context.TODO()

	testCases := []struct {
		like  *model.Like
		Valid bool
	}{
		{
			like:  model.NewLike(s.userID, s.articleID),
			Valid: true,
		},
		{
			like:  model.NewLike(0, s.articleID),
			Valid: false,
		},
		{
			like:  model.NewLike(s.userID, 0),
			Valid: false,
		},
	}

	for _, tc := range testCases {
		like, err := s.repo.Create(ctx, tc.like)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(like)
			s.likeID = like.ID

		} else if !tc.Valid {
			s.Error(err)
			s.Nil(like)
		}
	}
}

func (s *LikeTestSuite) TestB_GetLikeByID() {
	ctx := context.TODO()

	testCases := []struct {
		ID    uint
		Valid bool
	}{
		{
			ID:    s.likeID,
			Valid: true,
		},
		{
			ID:    0,
			Valid: false,
		},
	}

	for _, tc := range testCases {
		like, err := s.repo.GetByID(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(like)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(like)
		}
	}
}

func (s *LikeTestSuite) TestC_GetByUserID() {
	ctx := context.TODO()
	likes, err := s.repo.GetByUserID(ctx, s.userID)
	s.NoError(err)
	s.NotNil(likes)

}

func (s *LikeTestSuite) TestD_DeleteByID() {
	ctx := context.TODO()

	testCases := []struct {
		userID    uint
		ArticleID uint
		Valid     bool
	}{
		{
			userID:    s.userID,
			ArticleID: s.articleID,
			Valid:     true,
		},

		{
			userID:    0,
			ArticleID: s.articleID,
			Valid:     false,
		},

		{
			userID:    s.userID,
			ArticleID: 0,
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		err := s.repo.DeleteByArticleIDAndUserID(ctx, tc.ArticleID, tc.userID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}

	//Ensure user exists after deleting its like
	user, err := s.userRepo.GetUserByID(ctx, s.userID)
	s.Nil(err)
	s.NotNil(user)

	//Ensure article exists after deleting its like
	article, err := s.articleRepo.GetArticleByID(ctx, s.articleID)
	s.Nil(err)
	s.NotNil(article)
}

func TestLikeTestSuite(t *testing.T) {
	suite.Run(t, new(LikeTestSuite))
}
