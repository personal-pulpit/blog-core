package service

import (
	"blog/config"
	postgres_repository "blog/database/postgres/repo"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/internal/service/authentication"
	email "blog/pkg/email_manager"
	"blog/utils/hash"
	"context"
	"testing"

	auth_manager "blog/pkg/auth_manager"

	"github.com/stretchr/testify/suite"
)

type AuthTestSuite struct {
	suite.Suite
	authRepo repository.AuthPostgresRepository
	userRepo repository.UserPostgresRepository
	service  authentication.AuthService

	user               *model.User
	verificationCode   string
	password           string
	accessToken        string
	refreshToken       string
	resetPasswordToken string
}

func (s *AuthTestSuite) SetupSuite() {
	s.authRepo = postgres_repository.NewAuthPostgresRepository(db)
	s.userRepo = postgres_repository.NewUserPostgresRepository(db)

	authManager := auth_manager.NewAuthManger(redisCLI, config.GetConfigInstance().Jwt)

	hashManager := hash.NewHashManager(hash.DefaultHashParams)

	s.service = authentication.NewAuthenticateService(s.authRepo, s.userRepo, authManager, hashManager, email.NewEmailService(&config.GetConfigInstance().Email))
}

func (s *AuthTestSuite) TestA_Register() {
	ctx := context.TODO()

	testCases := []struct {
		firstName string
		lastName  string
		email     string
		biography string
		password  string
		Valid     bool
	}{
		{
			firstName: "first name1",
			lastName:  "last name1",
			email:     "example@gmail.com",
			biography: "valid user",
			password:  "password!!!",
			Valid:     true,
		},
		{
			firstName: "first name2",
			lastName:  "last name2",
			email:     "example@gmail.com",
			biography: "invalid user",
			password:  "!!!passwordsi",
			Valid:     false,
		},
	}

	for _, tc := range testCases {
		user, err := s.service.Register(ctx, tc.firstName, tc.lastName, tc.email, tc.biography, tc.password)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(user)

			s.user = user
			verificationCode := getVerificationCodeFromRedis(ctx)
			s.verificationCode = verificationCode
			s.password = tc.password
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(user)
		}
	}
}

func (s *AuthTestSuite) TestB_VerifyEmail() {
	ctx := context.TODO()

	testCases := []struct {
		userID uint
		opt    string
		Valid  bool
	}{
		{
			userID: s.user.ID,
			opt:    s.verificationCode,
			Valid:  true,
		},
		{
			userID: 14885684,
			opt:    s.verificationCode,
			Valid:  false,
		},
		{
			userID: s.user.ID,
			opt:    "fakefakefake",
			Valid:  false,
		},
	}

	for _, tc := range testCases {
		err := s.service.VerifyEmail(ctx, tc.opt, tc.userID)
		if tc.Valid {
			s.NoError(err)
		} else {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestC_Login() {
	ctx := context.TODO()

	testCases := []struct {
		email    string
		password string
		Valid    bool
	}{
		{
			email:    s.user.Email,
			password: s.password,
			Valid:    true,
		},
		{
			email:    "afakeemail@gmail.com",
			password: s.password,
			Valid:    false,
		},
		{
			email:    s.user.Email,
			password: "invalidPassword",
			Valid:    false,
		},
	}

	for _, tc := range testCases {
		user, accessToken, refreshToken, err := s.service.Login(ctx, tc.email, tc.password)
		if tc.Valid {
			s.NotNil(user)
			s.NotEmpty(accessToken)
			s.NotEmpty(refreshToken)
			s.NoError(err)

			s.accessToken = accessToken
			s.refreshToken = refreshToken
		} else {
			s.Error(err)
			s.Nil(user)
			s.Empty(accessToken)
			s.Empty(refreshToken)
		}
	}
}

func (s *AuthTestSuite) TestD_Authenticate() {
	ctx := context.TODO()

	testCases := []struct {
		accessToken string
		Valid       bool
	}{
		{
			accessToken: s.accessToken,
			Valid:       true,
		},
		{
			accessToken: "Invalid",
			Valid:       false,
		},
	}

	for _, tc := range testCases {
		user, err := s.service.Authenticate(ctx, tc.accessToken)
		if tc.Valid {
			s.NoError(err)
			s.NotNil(user)
			s.Equal(s.user.ID, user.ID)
			s.Equal(s.user.Email, user.Email)
			s.Equal(s.user.FirstName, user.FirstName)
			s.Equal(s.user.LastName, user.LastName)
			s.Equal(s.user.Biography, user.Biography)
		} else if !tc.Valid {
			s.Error(err)
			s.Nil(user)
		}
	}
}

func (s *AuthTestSuite) TestE_ChangePassword() {
	ctx := context.TODO()
	testCases := []struct {
		accessToken string
		oldPassword string
		newPassword string
		Valid       bool
	}{
		{
			accessToken: s.accessToken,
			oldPassword: s.password,
			newPassword: "newPassForMeVerySecure",
			Valid:       true,
		},
		{
			accessToken: "Invalid",
			oldPassword: s.password,
			newPassword: "newPassForMeVerySecure",
			Valid:       false,
		},
		{
			accessToken: s.accessToken,
			oldPassword: "Invalid",
			newPassword: "newPassForMeVerySecure",
			Valid:       false,
		},
	}

	for _, tc := range testCases {
		err := s.service.ChangePassword(ctx, tc.accessToken, tc.oldPassword, tc.newPassword)
		if tc.Valid {
			s.NoError(err)

			s.quickLogin(ctx, s.user.Email, tc.newPassword)

			s.password = tc.newPassword
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestF_RefreshToken() {
	ctx := context.TODO()

	testCases := []struct {
		refreshToken string
		accessToken  string
		Valid        bool
	}{
		{
			refreshToken: s.refreshToken,
			accessToken:  s.accessToken,
			Valid:        true,
		},
		{
			refreshToken: "invalid",
			accessToken:  s.accessToken,
			Valid:        false,
		},
		{
			refreshToken: s.refreshToken,
			accessToken:  "invalid",
			Valid:        false,
		},
	}

	for _, tc := range testCases {
		newAccessToken, err := s.service.RefreshToken(ctx, tc.refreshToken, tc.accessToken)
		if tc.Valid {
			s.NoError(err)
			s.NotEmpty(newAccessToken)

			user, err := s.service.Authenticate(ctx, s.accessToken)
			s.NoError(err)
			s.NotNil(user)
			s.Equal(s.user.ID, user.ID)
			s.Equal(s.user.Email, user.Email)
			s.Equal(s.user.FirstName, user.FirstName)
			s.Equal(s.user.LastName, user.LastName)
			s.Equal(s.user.Biography, user.Biography)

			s.accessToken = newAccessToken
		} else if !tc.Valid {
			s.Error(err)
			s.Empty(newAccessToken)
		}
	}
}

func (s *AuthTestSuite) TestG_DestroyRefreshToken() {
	ctx := context.TODO()

	err := s.service.DestroyRefreshToken(ctx, s.refreshToken)
	s.NoError(err)

	_, err = s.service.RefreshToken(ctx, s.refreshToken, s.accessToken)
	s.Error(err)
}

func (s *AuthTestSuite) TestH_SendResetPasswordVerification() {
	ctx := context.TODO()

	testCases := []struct {
		email string
		Valid bool
	}{
		{
			email: s.user.Email,
			Valid: true,
		},
		{
			email: "afakeemail@gmail.com",
			Valid: false,
		},
	}
	for _, tc := range testCases {
		resetPasswordToken, err := s.service.SendResetPasswordVerification(ctx, tc.email)
		if tc.Valid {
			s.NoError(err)
			s.resetPasswordToken = resetPasswordToken
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestI_SubmitResetPassword() {
	ctx := context.TODO()

	testCases := []struct {
		resetPasswordToken string
		newPassword        string
		Valid              bool
	}{
		{
			resetPasswordToken: s.resetPasswordToken,
			newPassword:        "newnewnewPassword",
			Valid:              true,
		},
		{
			resetPasswordToken: "invalid",
			newPassword:        "newPassForMeVerySecure",
			Valid:              false,
		},
	}

	for _, tc := range testCases {
		err := s.service.SubmitResetPassword(ctx, tc.resetPasswordToken, tc.newPassword)
		if tc.Valid {
			s.NoError(err)

			s.quickLogin(ctx, s.user.Email, tc.newPassword)

			s.password = tc.newPassword
		} else if !tc.Valid {
			s.Error(err)
		}
	}
}

func (s *AuthTestSuite) TestJ_DeleteAccount() {
	ctx := context.TODO()

	testCases := []struct {
		password    string
		userID  uint
		Valid       bool
	}{
		{
			password:    s.password,
			userID: s.user.ID,
			Valid:       true,
		},
		{
			password:    "invalid",
			userID:s.user.ID ,
			Valid:       false,
		},
		{
			password:    s.password,
			userID: 45465486416,
			Valid:       false,
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

func (s *AuthTestSuite) quickLogin(ctx context.Context, email string, password string) {
	user, _, _, err := s.service.Login(ctx, email, password)

	s.NoError(err)
	s.NotNil(user)
	s.user = user
}

func getVerificationCodeFromRedis(ctx context.Context) string {
	keys, err := redisCLI.Keys(ctx, "*").Result()
	if err != nil {
		panic(err)
	}

	codeKey := keys[0]

	code, err := redisCLI.Get(ctx, codeKey).Result()
	if err != nil {
		panic(err)
	}

	return code
}

func TestAuthSuite(t *testing.T) {
	t.Helper()

	suite.Run(t, new(AuthTestSuite))
}
