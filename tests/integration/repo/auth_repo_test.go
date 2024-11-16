package repo

import (
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	repo      repository.AuthPostgresRepository
	savedAuth *model.Auth
}

func (s *AuthTestSuite) SetupSuite() {
	s.repo = postgres_repository.NewAuthPostgresRepository(db)
}

func (s *AuthTestSuite) TestA_Create() {
	ctx := context.TODO()

	testCases := []struct {
		auth  *model.Auth
		Valid bool
	}{
		{
			auth:  model.NewAuth(1, "password",model.UserRole),
			Valid: true,
		},
	}

	for _, tc := range testCases {
		auth, err := s.repo.Create(ctx, tc.auth)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(auth)
			s.savedAuth = auth
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(auth)
		}
	}
}

func (s *AuthTestSuite) TestB_GetUserByID() {
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
			ID:    s.savedAuth.ID,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		auth, err := s.repo.GetUserAuth(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(auth)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(auth)
		}
	}
}

func (s *AuthTestSuite) TestC_ChangePassword() {
	ctx := context.TODO()

	testCases := []struct {
		ID          uint
		newPassword string
		Valid       bool
	}{
		{
			ID:          1000000000,
			newPassword: "newPassword",
			Valid:       false,
		},
		{
			ID:          s.savedAuth.ID,
			newPassword: "newPassword",
			Valid:       true,
		},
	}

	for _, tc := range testCases {
		err := s.repo.ChangePassword(ctx, tc.ID, tc.newPassword)
		if tc.Valid {
			s.NoError(err)
			s.savedAuth.HashedPassword = tc.newPassword
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestD_VerifyEmail() {
	ctx := context.TODO()

	testCases := []struct {
		ID    uint
		Valid bool
	}{
		{
			ID:    1000000,
			Valid: false,
		},
		{
			ID:    s.savedAuth.ID,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		err := s.repo.VerifyEmail(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestE_IncrementFailedLoginAttempts() {
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
			ID:    s.savedAuth.ID,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		err := s.repo.IncrementFailedLoginAttempts(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestF_ClearFailedLoginAttempts() {
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
			ID:    s.savedAuth.ID,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		err := s.repo.ClearFailedLoginAttempts(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestG_LockAccount() {
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
			ID:    s.savedAuth.ID,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		err := s.repo.LockAccount(ctx, tc.ID, 2*time.Second)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestH_UnlockAccount() {
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
			ID:    s.savedAuth.ID,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		err := s.repo.UnlockAccount(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}


func (s *AuthTestSuite) TestI_DeleteByID() {
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
			ID:    s.savedAuth.ID,
			Valid: true,
		},
	}

	for _, tc := range testCases {
		err := s.repo.DeleteByID(ctx, tc.ID)
		if tc.Valid {
			s.NoError(err)
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func TestAuthTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
