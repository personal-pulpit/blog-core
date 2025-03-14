package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	repo      repository.UserPostgresRepository
	savedUser *model.User
	commentRepo repository.CommentRepository
	articleRepo repository.ArticlePostgresRepository

	usersArticleID uint
	usersCommentID uint

}

func (s *UserTestSuite) SetupSuite() {
	//get main repo
	s.repo = postgres_repository.NewUserPostgresRepository(db)

	//get helper repos
	s.commentRepo = postgres_repository.NewCommentPostgresRepository(db)
	s.articleRepo = postgres_repository.NewArticlePostgresRepo(db)
}

func (s *UserTestSuite) TestA_Create() {
	ctx := context.TODO()

	testCases := []struct {
		User  *model.User
		Valid bool
	}{
		{
			User:  model.NewUser("firstName1", "lastName1", "example@gmail.com", "first user"),
			Valid: true,
		},
		{
			User:  model.NewUser("firstName1", "lastName1", "example@gmail.com", "first user"),
			Valid: false,
		},
	}

	for _, tc := range testCases {
		user, tx, err := s.repo.Create(ctx, tc.User)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(user)
			tx.Commit()

			s.savedUser = user

			article,err := s.articleRepo.Create(ctx, model.NewArticle("article1", "content", user.ID))
			s.Nil(err)
			s.NotNil(article)
			
			s.usersArticleID = article.ID

			comment,err := s.commentRepo.Create(ctx, model.NewComment(user.ID, article.ID, "comment1"))
			s.Nil(err)
			s.NotNil(comment)

			s.usersCommentID = comment.ID
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(user)
		}
	}
}

func (s *UserTestSuite) TestB_GetUserByID() {
	ctx := context.TODO()

	testCases := []struct {
		ID    uint
		Valid bool
	}{
		{
			ID:    1000,
			Valid: false,
		},
		{
			ID:    s.savedUser.ID,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		user, err := s.repo.GetUserByID(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(user)
			s.NotNil(user.Comments)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(user)
		}
	}
}

func (s *UserTestSuite) TestC_CreateWithRollback() {
	ctx := context.TODO()
	user := model.NewUser("firstName2", "lastName2", "example2@gmail.com", "fist user")

	savedUser, tx, err := s.repo.Create(ctx, user)
	s.NoError(err)
	s.NotNil(savedUser)

	nilUser, err := s.repo.GetUserByID(ctx, savedUser.ID)
	s.Error(err)
	s.Nil(nilUser)

	tx.Rollback()

	nilUser, err = s.repo.GetUserByID(ctx, savedUser.ID)
	s.Error(err)
	s.Nil(nilUser)
}

func (s *UserTestSuite) TestD_GetUserByEmail() {
	ctx := context.TODO()

	testCases := []struct {
		email string
		Valid bool
	}{
		{
			email: "invalid",
			Valid: false,
		},
		{
			email: s.savedUser.Email,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		user, err := s.repo.GetUserByEmail(ctx, tc.email)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(user)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(user)
		}
	}
}

func (s *UserTestSuite) TestE_ModifyUser() {
	ctx := context.TODO()

	testCases := []struct {
		ID        uint
		firstName string
		lastName  string
		biography string
		Valid     bool
	}{
		{
			ID:        1000,
			firstName: "valid",
			lastName:  "valid",
			biography: "valid",
			Valid:     false,
		},
		{
			ID:        s.savedUser.ID,
			firstName: "valid",
			lastName:  "valid",
			biography: "valid",
			Valid:     true,
		},
	}

	for _, tc := range testCases {
		user, err := s.repo.UpdateByID(ctx, tc.ID, tc.firstName, tc.lastName, tc.biography)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(user)
			s.savedUser = user
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(user)
		}
	}
}

func (s *UserTestSuite) TestF_DeleteUser() {
	ctx := context.TODO()	
	testCases := []struct {
		userID uint
		Valid  bool
	}{
		{
			userID: s.savedUser.ID,
			Valid:  true,
		},
		{
			userID: 1000,
			Valid:  false,
		},
	}

	// Ensure the user exists before deleting
	_, err := s.repo.GetUserByID(ctx, s.savedUser.ID)
	s.NoError(err)

	for _, tc := range testCases {
		err := s.repo.DeleteByID(ctx, tc.userID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}

	// Ensure the user is deleted 
	_, err = s.repo.GetUserByID(ctx, s.savedUser.ID)
	s.Error(err)

	// Ensure the user's article is deleted 
	_,err = s.articleRepo.GetArticleByID(ctx,s.usersArticleID)
	s.Error(err)

	// Ensure the user's comment is deleted 
	_,err = s.commentRepo.GetByID(ctx,s.usersCommentID)
	s.Error(err)
}
func TestUserTestSuite(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}
