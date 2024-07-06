package authentication

import (
	"blog/internal/model"
	"blog/internal/repository"
	"blog/pkg/auth_manager"
	email "blog/pkg/email_manager"
	"blog/utils/hash"
	"blog/utils/random"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	OTPExpr                    = time.Minute * 10    //120 second
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
	VerifyEmail(otp string, userID model.ID) error
	SendResetPasswordVerification(email string) (string, time.Duration, error)
	SubmitResetPassword(token string, newPassword string) error
	ChangePassword(accessToken string, oldPassword string, newPassword string) error
	Authenticate(accessToken string) (*model.User, error)
	RefreshToken(refreshToken string, accessToken string) (string, error)
	DeleteAccount(ID model.ID, password string) error
	Logout(token string) error
}
type authenticateManager struct {
	uniqueId         string
	userPostgresRepo repository.UserPostgresRepository
	authPostgresRepo repository.AuthPostgresRepository
	authManager      auth_manager.AuthManager
	hashManager      *hash.HashManager
	emailService     email.EmailService
}

func NewAuthenticateService(authPostgresRepo repository.AuthPostgresRepository, userPostgresRepo repository.UserPostgresRepository, authManager auth_manager.AuthManager, hashManager *hash.HashManager, emailService email.EmailService) AuthService {
	return &authenticateManager{
		authPostgresRepo: authPostgresRepo,
		userPostgresRepo: userPostgresRepo,
		authManager:      authManager,
		hashManager:      hashManager,
		emailService:     emailService,
	}
}

func (a *authenticateManager) Register(firstName, lastName, email, biography, password string) (*model.User, string, error) {
	userModel := model.NewUser(firstName, lastName, email, biography, model.UserRole)

	savedUser, tx, err := a.userPostgresRepo.Create(userModel)

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

	_, err = a.authPostgresRepo.Create(authModel)

	if err != nil {
		tx.Rollback()

		return nil, "", ErrCreateAuthStore
	}

	a.uniqueId = fmt.Sprintf("%d", random.GenerateUniqueId())

	verifyEmailToken, err := a.authManager.GenerateToken(
		auth_manager.VerifyEmail,
		auth_manager.NewTokenClaims(savedUser.ID, model.UserRole, auth_manager.VerifyEmail),
		VerifyEmailTokenExpr,
	)

	if err != nil {
		tx.Rollback()

		return nil, "", ErrCreateEmailToken
	}

	otp, err := a.authManager.SetOTP(a.uniqueId, OTPExpr)

	if err != nil {
		tx.Rollback()

		return nil, "", err
	}

	err = a.emailService.SendVerificationEmail(userModel.Email, otp)

	if err != nil {
		tx.Rollback()

		return nil, "", err
	}

	return savedUser, verifyEmailToken, nil
}

func (a *authenticateManager) VerifyEmail(otp string, userId model.ID) error {
	savedOTP, err := a.authManager.GetOTP(a.uniqueId)

	if err != nil {
		return err
	}

	if savedOTP != otp {
		return ErrVerifyEmail
	}

	err = a.authPostgresRepo.VerifyEmail(userId)
	if err != nil {
		return ErrVerifyEmail
	}

	return nil
}
func (a *authenticateManager) Login(email string, password string) (*model.User, string, string, error) {
	userModel, err := a.userPostgresRepo.GetUserByEmail(email)
	if errors.Is(err, repository.ErrUniqueConstraint) {
		return nil, "", "", repository.ErrUniqueConstraint
	} else if err != nil {
		return nil, "", "", ErrCreateUser
	}
	auth, err := a.authPostgresRepo.GetUserAuth(userModel.ID)
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
			err = a.authPostgresRepo.UnlockAccount(auth.ID)
			if err != nil {
				return nil, "", "", ErrUnlockAccount
			}

			err = a.authPostgresRepo.ClearFailedLoginAttempts(auth.ID)
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
		err = a.authPostgresRepo.LockAccount(auth.ID, LockAccountDuration)
		if err != nil {
			return nil, "", "", ErrLockAccount
		}
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.HashedPassword)
	if !validPassword {
		err = a.authPostgresRepo.IncrementFailedLoginAttempts(userModel.ID)
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

	err = a.authPostgresRepo.ClearFailedLoginAttempts(auth.ID)
	if err != nil {
		return nil, "", "", ErrClearFailedLoginAttempts
	}

	err = a.emailService.SendWelcomeEmail(email, userModel.FirstName)
	if err != nil {
		return nil, "", "", err
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

	user, err := a.userPostgresRepo.GetUserByID(tokenClaims.ID)
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

	auth, err := a.authPostgresRepo.GetUserAuth(user.ID)
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

	err = a.authPostgresRepo.ChangePassword(user.ID, newPasswordHash)
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

	_, err = a.authPostgresRepo.GetUserAuth(rftClaims.ID)
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
	user, err := a.userPostgresRepo.GetUserByEmail(email)
	if err != nil {
		return "", 0, err
	}

	auth, err := a.authPostgresRepo.GetUserAuth(user.ID)
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

	err = a.emailService.SendResetPasswordEmail(email, "example.com", user.FirstName, "10")
	if err != nil {
		return "", 0, err
	}

	return resetPasswordToken, ResetPasswordTokenExpr, nil
}

func (a *authenticateManager) SubmitResetPassword(token string, newPassword string) error {
	tokenClaims, err := a.authManager.DecodeToken(token, auth_manager.ResetPassword)
	if err != nil {
		return ErrAccessDenied
	}

	auth, err := a.authPostgresRepo.GetUserAuth(tokenClaims.ID)
	if err != nil {
		return ErrAccessDenied
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return ErrHashingPassword
	}

	err = a.authPostgresRepo.ChangePassword(auth.ID, newPasswordHash)
	if err != nil {
		return ErrChangePassword
	}

	return nil
}

func (a *authenticateManager) DeleteAccount(ID model.ID, password string) error {
	auth, err := a.authPostgresRepo.GetUserAuth(ID)
	if err != nil {
		return ErrNotFound
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.HashedPassword)
	if !validPassword {
		return ErrDeleteUser
	}

	err = a.authPostgresRepo.DeleteByID(ID)
	if err != nil {
		return ErrDeleteUser
	}

	err = a.userPostgresRepo.DeleteByID(ID)
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
