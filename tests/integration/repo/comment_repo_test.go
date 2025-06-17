package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type CommentTestSuite struct {
	suite.Suite
	repo        repository.CommentRepository
	userRepo    repository.UserRepository
	articleRepo repository.ArticleRepository
	commentID   uint
	articleID   uint
	userID      uint
}

func (s *CommentTestSuite) SetupSuite() {
	s.repo = postgres_repository.NewCommentPostgresRepository(db)
	s.userRepo = postgres_repository.NewUserPostgresRepository(db)
	s.articleRepo = postgres_repository.NewArticlePostgresRepo(db)
	categoryRepo := postgres_repository.NewCategoryRepository(db)

	user, tx, err := s.userRepo.Create(context.TODO(), model.NewUser("user2", "user2", "<TEST-COMMENT>", "user2"))
	s.Nil(err)
	s.NotNil(user)
	tx.Commit()

	s.userID = user.ID

	category, err := categoryRepo.CreateCategory(context.TODO(), model.NewCategory("comment-category1"))
	s.Nil(err)
	s.NotNil(category)

	article, err := s.articleRepo.Create(context.TODO(), model.NewArticle("article1", "content", s.userID, []model.Category{*category}))
	s.Nil(err)
	s.NotNil(article)

	s.articleID = article.ID
}

func (s *CommentTestSuite) TestA_CreateComment() {
	ctx := context.TODO()

	testCases := []struct {
		comment *model.Comment
		Valid   bool
	}{
		{
			comment: model.NewComment(s.userID, s.articleID, "comment1"),
			Valid:   true,
		},
		{
			comment: model.NewComment(0, s.articleID, "comment1"),
			Valid:   false,
		},
		{
			comment: model.NewComment(s.userID, 0, "comment1"),
			Valid:   false,
		},
		// {
		// 	comment: model.NewComment(s.userID, s.articleID, ""),
		// 	Valid:   false,
		// },
	}

	for _, tc := range testCases {
		comment, err := s.repo.Create(ctx, tc.comment)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(comment)
			s.commentID = comment.ID

		} else if !tc.Valid {
			s.Error(err)
			s.Nil(comment)
		}
	}
}

func (s *CommentTestSuite) TestB_GetAllComments() {
	ctx := context.TODO()

	comments, err := s.repo.GetAll(ctx)
	s.NoError(err)
	s.NotNil(comments)
}

func (s *CommentTestSuite) TestC_GetCommentByID() {
	ctx := context.TODO()

	testCases := []struct {
		ID    uint
		Valid bool
	}{
		{
			ID:    s.commentID,
			Valid: true,
		},
		{
			ID:    0,
			Valid: false,
		},
	}

	for _, tc := range testCases {
		comment, err := s.repo.GetByID(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(comment)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(comment)
		}
	}
}

func (s *CommentTestSuite) TestD_UpdateComment() {
	ctx := context.TODO()

	testCases := []struct {
		commentID uint
		content   string
		Valid     bool
	}{
		{
			commentID: s.commentID,
			content:   "comment 0",
			Valid:     true,
		},
		{
			commentID: 0,
			content:   "comment 0",
			Valid:     false,
		},
		// {
		// 	commentID: s.commentID,
		// 	content:   "",
		// 	Valid:     false,
		// },
	}

	for _, tc := range testCases {
		savedComment, err := s.repo.UpdateByID(ctx, tc.commentID, tc.content)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(savedComment)
			s.NotEmpty(savedComment.UpdatedAt)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *CommentTestSuite) TestE_DeleteCommentByID() {
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
		err := s.repo.DeleteByID(ctx, tc.commentID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}

	//Ensure user exists after deleting its comment
	user, err := s.userRepo.GetUserByID(ctx, s.userID)
	s.Nil(err)
	s.NotNil(user)

	//Ensure article exists after deleting its comment
	article, err := s.articleRepo.GetArticleByID(ctx, s.articleID)
	s.Nil(err)
	s.NotNil(article)
}

func TestCommentTestSuite(t *testing.T) {
	suite.Run(t, new(CommentTestSuite))
}
