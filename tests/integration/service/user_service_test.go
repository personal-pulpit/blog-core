package service

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/internal/service/user"
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserTestSuite struct {
	suite.Suite
	userRepo repository.UserPostgresRepository

	service user.UserService

	user     *model.User
	password string
}

func (s *UserTestSuite) SetupSuite() {
	s.userRepo = postgres_repository.NewUserPostgresRepository(db)
	authRepo := postgres_repository.NewAuthPostgresRepository(db)

	s.service = user.NewUserService(s.userRepo, authRepo)

	userModel := model.NewUser("user1 firstName", "user1 lastName", "afakeonce@fake.come", "user1 biography")
	user, tx, err := s.userRepo.Create(context.TODO(), userModel)
	s.NoError(err)
	tx.Commit()

	password := "password123456"

	userAuth := model.NewAuth(user.ID, password, model.UserRole)
	_, err = authRepo.Create(context.TODO(), userAuth)
	s.NoError(err)

	err = authRepo.VerifyEmail(context.TODO(), user.ID)
	s.NoError(err)

	s.user = user
	s.password = password
}

func (s *UserTestSuite) TestA_GetUser() {
	ctx := context.TODO()

	testCases := []struct {
		ID    uint
		Valid bool
	}{
		{
			ID:    s.user.ID,
			Valid: true,
		},
		{
			ID:    8880008,
			Valid: false,
		},
	}

	for _, tc := range testCases {
		user, err := s.service.GetUserProfile(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(user)

		} else if !tc.Valid {
			s.Error(err)
			s.Nil(user)
		}
	}
}

func (s *UserTestSuite) TestB_UpdateUser() {
	ctx := context.TODO()

	testCases := []struct {
		userID    uint
		firstName string
		lastName  string
		biography string
		Valid     bool
	}{
		{
			userID:    s.user.ID,
			firstName: "first name1",
			lastName:  "last name1",
			biography: "valid user",
			Valid:     true,
		},
		{
			userID:    14885684,
			firstName: "first name1",
			lastName:  "last name1",
			biography: "invalid user",
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		user, err := s.service.UpdateProfile(ctx, tc.userID, tc.firstName, tc.lastName, tc.biography)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(user)
			s.Equal(s.user.ID, user.ID)
			s.NotEqual(s.user.FirstName, user.FirstName)
			s.NotEqual(s.user.LastName, user.LastName)
			s.NotEqual(s.user.Biography, user.Biography)

			s.user = user
		} else {
			s.Error(err)
			s.Nil(user)

		}
	}
}

func (s *UserTestSuite) TestC_DeleteAccount() {
	ctx := context.TODO()

	testCases := []struct {
		password string
		userID   uint
		Valid    bool
	}{
		{
			password: s.password,
			userID:   s.user.ID,
			Valid:    true,
		},
		{
			password: "invalid",
			userID:   s.user.ID,
			Valid:    false,
		},
		{
			password: s.password,
			userID:   45465486416,
			Valid:    false,
		},
	}

	for _, tc := range testCases {
		err := s.service.DeleteAccount(ctx, tc.userID, tc.password)
		if tc.Valid {
			s.NoError(err)

		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func TestUserSuite(t *testing.T) {
	t.Helper()

	suite.Run(t, new(UserTestSuite))
}
