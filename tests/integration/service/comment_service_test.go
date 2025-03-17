package service

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/internal/service/comment"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommentTestSuite struct {
	suite.Suite
	commentRepo repository.CommentRepository

	service comment.CommentService

	commentID uint
	articleID uint
	userID    uint
}

func (s *CommentTestSuite) SetupSuite() {
	s.commentRepo = postgres_repository.NewCommentPostgresRepository(db)
	userRepo := postgres_repository.NewUserPostgresRepository(db)
	articleRepo := postgres_repository.NewArticlePostgresRepo(db)

	userModel := model.NewUser("user1 firstName", "user1 lastName", "afakeonce@fake.come", "user1 biography")
	user, tx, err := userRepo.Create(context.TODO(), userModel)
	s.NoError(err)
	s.NotNil(user)
	tx.Commit()

	s.userID = user.ID

	article, err := articleRepo.Create(context.TODO(), model.NewArticle("article1", "content of article1", user.ID))
	s.NoError(err)
	s.NotNil(article)

	s.articleID = article.ID

	s.service = comment.NewCommentService(s.commentRepo)
}

func (s *CommentTestSuite) TestA_CreateComment() {
	ctx := context.TODO()

	comment, err := s.service.AddComment(ctx, "comment1", s.userID, s.articleID)

	s.NoError(err)
	s.NotNil(comment)

	s.commentID = comment.ID
}

func (s *CommentTestSuite) TestB_UpdateComment() {
	ctx := context.TODO()

	testCases := []struct {
		commentID uint
		content   string
		Valid     bool
	}{
		{
			commentID: s.commentID,
			content:   "changed comment1",
			Valid:     true,
		},
		{
			commentID: 0,
			content:   "failed  content1",
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		comment, err := s.service.UpdateComment(ctx, tc.commentID, tc.content)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(comment)
			s.Equal(s.commentID, comment.ID)
		} else {
			s.Error(err)
			s.Nil(comment)
		}
	}
}

func (s *CommentTestSuite) TestC_DeleteComment() {
	ctx := context.TODO()

	testCases := []struct {
		commentID uint
		Valid     bool
	}{
		{
			commentID: s.commentID,
			Valid:     true,
		},
		{
			commentID: 0,
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		err := s.service.DeleteComment(ctx, tc.commentID)
		if tc.Valid {
			s.NoError(err)

		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func TestCommentSuite(t *testing.T) {
	t.Helper()

	suite.Run(t, new(CommentTestSuite))
}
