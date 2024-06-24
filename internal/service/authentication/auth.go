package authentication

import (
	"blog/internal/model"
	"blog/internal/repository"
	"blog/pkg/auth_manager"
	"blog/utils/hash"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	ResetPasswordTokenExpr     = time.Minute * 10    // 10 minutes
	VerifyEmailTokenExpr       = time.Minute * 5     // 5 minutes
	AccessTokenExpr            = time.Hour * 24 * 2  // 2 days
	RefreshTokenExpr           = time.Hour * 24 * 14 // 2 weeks
	LockAccountDuration        = time.Second * 5
	MaximumFailedLoginAttempts = 5
)

type AuthService interface {
	Register(FirstName, lastName, email, biography, password string) (*model.User, string, error)
	Login(email string, password string) (*model.User, string, string, error)
	VerifyEmail(verifyEmailToken string) error
	SendResetPasswordVerification(email string) (string, time.Duration, error)
	SubmitResetPassword(token string, newPassword string) error
	ChangePassword(accessToken string, oldPassword string, newPassword string) error
	Authenticate(accessToken string) (*model.User, error)
	RefreshToken(refreshToken string, accessToken string) (string, error)
	DeleteAccount(ID model.ID, password string) error
	Logout(token string)error
}
type authenticateManager struct {
	userMysqlRepo repository.UserMysqlRepository
	authMysqlRepo repository.AuthMysqlRepository
	authManager   auth_manager.AuthManager
	hashManager   *hash.HashManager
}

func NewAuthenticateService(authMysqlRepo repository.AuthMysqlRepository, userMysqlRepo repository.UserMysqlRepository, authManager auth_manager.AuthManager, hashManager *hash.HashManager) AuthService {
	return &authenticateManager{
		authMysqlRepo: authMysqlRepo,
		userMysqlRepo: userMysqlRepo,
		authManager:   authManager,
		hashManager:   hashManager,
	}
}

func (a *authenticateManager) Register(FirstName, lastName, email, biography, password string) (*model.User, string, error) {
	userModel := model.NewUser(FirstName, lastName, email, biography, model.UserRole)
	savedUser, tx, err := a.userMysqlRepo.Create(userModel)
	if errors.Is(err, repository.ErrUniqueConstraint) {
		return nil, "", repository.ErrUniqueConstraint
	} else if err != nil {
		return nil, "", ErrCreateUser
	}
	passwordHash, err := a.hashManager.HashPassword(password)
	if err != nil {
		tx.Rollback()
		return nil, "", ErrHashingPassword
	}

	authModel := model.NewAuth(savedUser.ID, passwordHash)
	_, err = a.authMysqlRepo.Create(authModel)
	if err != nil {
		return nil, "", ErrCreateAuthStore
	}

	verifyEmailToken, err := a.authManager.GenerateToken(
		auth_manager.VerifyEmail,
		auth_manager.NewTokenClaims(savedUser.ID, userModel.Role, auth_manager.VerifyEmail),
		VerifyEmailTokenExpr,
	)
	if err != nil {
		return nil, "", ErrCreateEmailToken
	}

	return savedUser, verifyEmailToken, nil
}
func (a *authenticateManager) VerifyEmail(verifyEmailToken string) error {
	tokenClaims, err := a.authManager.DecodeToken(verifyEmailToken, auth_manager.VerifyEmail)
	if err != nil {
		return ErrAccessDenied
	}

	err = a.authMysqlRepo.VerifyEmail(tokenClaims.ID)
	if err != nil {
		return ErrVerifyEmail
	}

	err = a.authManager.Destroy(verifyEmailToken)
	if err != nil {
		return ErrDestroyToken
	}

	return nil
}
func (a *authenticateManager) Login(email string, password string) (*model.User, string, string, error) {
	userModel, err := a.userMysqlRepo.GetUserByEmail(email)
	if errors.Is(err, repository.ErrUniqueConstraint) {
		return nil, "", "", repository.ErrUniqueConstraint
	} else if err != nil {
		return nil, "", "", ErrCreateUser
	}
	auth, err := a.authMysqlRepo.GetUserAuth(userModel.ID)
	if err != nil {
		return nil, "", "", ErrInvalidEmailOrPassword
	}

	if !auth.EmailVerified {
		return nil, "", "", ErrEmailNotVerified
	}

	if auth.AccountLockedUntil != 0 {
		now := time.Now()
		lockTime := time.Unix(auth.AccountLockedUntil, 0)
		if now.After(lockTime) {
			err = a.authMysqlRepo.UnlockAccount(auth.ID)
			if err != nil {
				return nil, "", "", ErrUnlockAccount
			}

			err = a.authMysqlRepo.ClearFailedLoginAttempts(auth.ID)
			if err != nil {
				return nil, "", "", ErrClearFailedLoginAttempts
			}

			auth.AccountLockedUntil = 0
		}
	}

	if auth.AccountLockedUntil != 0 {
		lockTime := time.Unix(auth.AccountLockedUntil, 0)
		return nil, "", "", fmt.Errorf("%w until %v", ErrAccountLocked, lockTime)
	}

	if auth.FailedLoginAttempts+1 == MaximumFailedLoginAttempts {
		err = a.authMysqlRepo.LockAccount(auth.ID, LockAccountDuration)
		if err != nil {
			return nil, "", "", ErrLockAccount
		}
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.HashedPassword)
	if !validPassword {
		err = a.authMysqlRepo.IncrementFailedLoginAttempts(userModel.ID)
		if err != nil {
			return nil, "", "", ErrInvalidEmailOrPassword
		}

		return nil, "", "", ErrInvalidEmailOrPassword
	}

	accessToken, err := a.authManager.GenerateToken(auth_manager.AccessToken, auth_manager.NewTokenClaims(userModel.ID, userModel.Role, auth_manager.AccessToken), AccessTokenExpr)
	if err != nil {
		return nil, "", "", ErrGenerateToken
	}

	refreshToken, err := a.authManager.GenerateToken(auth_manager.RefreshToken, auth_manager.NewTokenClaims(userModel.ID, userModel.Role, auth_manager.RefreshToken), RefreshTokenExpr)
	if err != nil {
		return nil, "", "", ErrGenerateToken
	}

	err = a.authMysqlRepo.ClearFailedLoginAttempts(auth.ID)
	if err != nil {
		return nil, "", "", ErrClearFailedLoginAttempts
	}

	return userModel, accessToken, refreshToken, nil
}
func (a *authenticateManager) Authenticate(accessToken string) (*model.User, error) {

	tokenClaims, err := a.authManager.DecodeToken(accessToken, auth_manager.AccessToken)
	if err != nil {
		return nil, ErrAccessDenied
	}

	if len(strings.TrimSpace(string(tokenClaims.ID))) == 0 {
		return nil, ErrAccessDenied
	}

	user, err := a.userMysqlRepo.GetUserByID(tokenClaims.ID)
	if err != nil {
		return nil, ErrAccessDenied
	}

	return user, nil
}
func (a *authenticateManager) ChangePassword(accessToken string, oldPassword string, newPassword string) error {
	user, err := a.Authenticate(accessToken)
	if err != nil {
		return err
	}

	auth, err := a.authMysqlRepo.GetUserAuth(user.ID)
	if err != nil {
		return ErrNotFound
	}

	validPassword := a.hashManager.CheckPasswordHash(oldPassword, auth.HashedPassword)
	if !validPassword {
		return ErrInvalidPassword
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return ErrHashingPassword
	}

	err = a.authMysqlRepo.ChangePassword(user.ID, newPasswordHash)
	if err != nil {
		return ErrChangePassword
	}

	return nil
}

func (a *authenticateManager) RefreshToken(refreshToken string, accessToken string) (string, error) {
	rftClaims, err := a.authManager.DecodeToken(refreshToken, auth_manager.RefreshToken)
	if err != nil {
		return "", ErrAccessDenied
	}

	_, err = a.authManager.DecodeToken(accessToken, auth_manager.AccessToken)
	if err != nil {
		return "", ErrAccessDenied
	}

	_, err = a.authMysqlRepo.GetUserAuth(rftClaims.ID)
	if err != nil {
		return "", ErrAccessDenied
	}

	newAccessToken, err := a.authManager.GenerateToken(auth_manager.AccessToken, auth_manager.NewTokenClaims(rftClaims.ID, rftClaims.Role, auth_manager.AccessToken), AccessTokenExpr)
	if err != nil {
		return "", ErrGenerateToken
	}

	err = a.authManager.Destroy(accessToken)
	if err != nil {
		return "", ErrDestroyToken
	}

	return newAccessToken, nil
}

func (a *authenticateManager) SendResetPasswordVerification(email string) (token string, timeout time.Duration, _ error) {
	user, err := a.userMysqlRepo.GetUserByEmail(email)
	if err != nil {
		return "", 0, err
	}

	auth, err := a.authMysqlRepo.GetUserAuth(user.ID)
	if err != nil {
		return "", 0, err
	}

	if !auth.EmailVerified {
		return "", 0, ErrEmailNotVerified
	}

	if auth.FailedLoginAttempts >= MaximumFailedLoginAttempts {
		return "", 0, fmt.Errorf("%w until: %v", ErrAccountLocked, auth.AccountLockedUntil)
	}

	resetPasswordToken, err := a.authManager.GenerateToken(auth_manager.ResetPassword, auth_manager.NewTokenClaims(auth.ID, user.Role, auth_manager.ResetPassword), ResetPasswordTokenExpr)
	if err != nil {
		return "", 0, ErrGenerateToken
	}

	return resetPasswordToken, ResetPasswordTokenExpr, nil
}

func (a *authenticateManager) SubmitResetPassword(token string, newPassword string) error {
	tokenClaims, err := a.authManager.DecodeToken(token, auth_manager.ResetPassword)
	if err != nil {
		return ErrAccessDenied
	}

	auth, err := a.authMysqlRepo.GetUserAuth(tokenClaims.ID)
	if err != nil {
		return ErrAccessDenied
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return ErrHashingPassword
	}

	err = a.authMysqlRepo.ChangePassword(auth.ID, newPasswordHash)
	if err != nil {
		return ErrChangePassword
	}

	return nil
}

func (a *authenticateManager) DeleteAccount(ID model.ID, password string) error {
	auth, err := a.authMysqlRepo.GetUserAuth(ID)
	if err != nil {
		return ErrNotFound
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.HashedPassword)
	if !validPassword {
		return ErrDeleteUser
	}

	err = a.authMysqlRepo.DeleteByID(ID)
	if err != nil {
		return ErrDeleteUser
	}

	err = a.userMysqlRepo.DeleteByID(ID)
	if err != nil {
		return ErrDeleteUser
	}

	return nil
}
func (a *authenticateManager) Logout(token string) error {
	err := a.authManager.Destroy(token)
	if err != nil {
		return ErrNotFound
	}
	return nil
}