package authentication

import (
	"blog/internal/model"
	"blog/internal/repository"
	"blog/pkg/auth_manager"
	email "blog/pkg/email_manager"
	"blog/utils/hash"
	"blog/utils/random"
	"context"
	"errors"
	"fmt"
	"strconv"
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
	Register(ctx context.Context, firstName, lastName, email, biography, password string) (*model.User, error)
	Login(ctx context.Context, email string, password string) (*model.User, string, string, error)
	VerifyEmail(ctx context.Context, otp string, userID uint) error
	SendResetPasswordVerification(ctx context.Context, email string) (string, error)
	SubmitResetPassword(ctx context.Context, token string, newPassword string) error
	ChangePassword(ctx context.Context, accessToken string, oldPassword string, newPassword string) error
	Authenticate(ctx context.Context, accessToken string) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error)
	SetAccessTokenIntoBlacklist(ctx context.Context, accessToken string, accessTokenExpr time.Duration) error
	IsAccessTokenBlacklisted(ctx context.Context, accessToken string) bool
	DeleteAccount(ctx context.Context,userID uint, password string) error
	DestroyRefreshToken(ctx context.Context, token string) error
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

func (a *authenticateManager) Register(ctx context.Context, firstName, lastName, email, biography, password string) (*model.User, error) {
	userModel := model.NewUser(firstName, lastName, email, biography)

	savedUser, tx, err := a.userPostgresRepo.Create(ctx, userModel)

	if errors.Is(err, repository.ErrUniqueConstraint) {
		return nil, repository.ErrUniqueConstraint

	} else if err != nil {
		return nil, ErrCreateUser
	}

	passwordHash, err := a.hashManager.HashPassword(password)

	if err != nil {
		tx.Rollback()
		return nil, ErrHashingPassword
	}

	authModel := model.NewAuth(savedUser.ID, passwordHash, model.UserRole)

	_, err = a.authPostgresRepo.Create(ctx, authModel)

	if err != nil {
		tx.Rollback()

		return nil, ErrCreateAuthStore
	}

	a.uniqueId = fmt.Sprintf("%d", random.GenerateUniqueId())

	emailVerificationCode, err := a.authManager.GenerateVerificationCode(ctx, a.uniqueId)

	if err != nil {
		tx.Rollback()

		return nil, ErrGenerateVerificationCode
	}

	err = a.emailService.SendVerificationEmail(userModel.Email, emailVerificationCode)

	if err != nil {
		tx.Rollback()

		return nil, err
	}

	tx.Commit()

	return savedUser, nil
}

func (a *authenticateManager) VerifyEmail(ctx context.Context, verificationCode string, userID uint) error {
	isValid, err := a.authManager.CompareVerificationCode(ctx, a.uniqueId, verificationCode)

	if err != nil {
		return err
	}

	if !isValid {
		return ErrVerifyEmail
	}

	err = a.authPostgresRepo.VerifyEmail(ctx, userID)
	if err != nil {
		return ErrVerifyEmail
	}

	return nil
}
func (a *authenticateManager) Login(ctx context.Context, email string, password string) (*model.User, string, string, error) {
	userModel, err := a.userPostgresRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", ErrInvalidEmailOrPassword
	}

	auth, err := a.authPostgresRepo.GetUserAuth(ctx, userModel.ID)
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
			err = a.authPostgresRepo.UnlockAccount(ctx, auth.ID)
			if err != nil {
				return nil, "", "", ErrUnlockAccount
			}

			err = a.authPostgresRepo.ClearFailedLoginAttempts(ctx, auth.ID)
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
		err = a.authPostgresRepo.LockAccount(ctx, auth.ID, LockAccountDuration)
		if err != nil {
			return nil, "", "", ErrLockAccount
		}
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.HashedPassword)
	if !validPassword {
		err = a.authPostgresRepo.IncrementFailedLoginAttempts(ctx, userModel.ID)
		if err != nil {
			return nil, "", "", ErrInvalidEmailOrPassword
		}

		return nil, "", "", ErrInvalidEmailOrPassword
	}

	accessToken, err := a.authManager.GenerateAccessToken(ctx, userModel.ID, auth.Role)
	if err != nil {
		return nil, "", "", ErrGenerateToken
	}

	refreshToken, err := a.authManager.GenerateRefreshToken(ctx, userModel.ID, "not implemented", "not implemented")
	if err != nil {
		return nil, "", "", ErrGenerateToken
	}

	err = a.authPostgresRepo.ClearFailedLoginAttempts(ctx, auth.ID)
	if err != nil {
		return nil, "", "", ErrClearFailedLoginAttempts
	}

	err = a.emailService.SendWelcomeEmail(email, userModel.FirstName)
	if err != nil {
		return nil, "", "", err
	}

	return userModel, accessToken, refreshToken, nil
}
func (a *authenticateManager) Authenticate(ctx context.Context, accessToken string) (*model.User, error) {

	tokenClaims, err := a.authManager.DecodeAccessToken(ctx, accessToken)
	if err != nil {
		return nil, ErrAccessDenied
	}

	if len(strings.TrimSpace(string(tokenClaims.UserID))) == 0 || len(strings.TrimSpace(tokenClaims.Role)) == 0 {
		return nil, ErrAccessDenied
	} else if !(tokenClaims.Role == fmt.Sprintf("%d", model.UserRole)) && !(tokenClaims.Role == fmt.Sprintf("%d", model.AdminRole)) {
		return nil, ErrAccessDenied
	}

	ID := convertString2Uint(tokenClaims.UserID)

	user, err := a.userPostgresRepo.GetUserByID(ctx, ID)
	if err != nil {
		return nil, ErrAccessDenied
	}

	return user, nil
}
func (a *authenticateManager) ChangePassword(ctx context.Context, accessToken string, oldPassword string, newPassword string) error {
	user, err := a.Authenticate(ctx, accessToken)
	if err != nil {
		return err
	}

	auth, err := a.authPostgresRepo.GetUserAuth(ctx, user.ID)
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

	err = a.authPostgresRepo.ChangePassword(ctx, user.ID, newPasswordHash)
	if err != nil {
		return ErrChangePassword
	}

	return nil
}

func (a *authenticateManager) RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error) {
	refreshTokenClaims, err := a.authManager.DecodeRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", ErrAccessDenied
	}

	_, err = a.authManager.DecodeAccessToken(ctx, accessToken)

	if err != nil {
		if errors.Is(err, auth_manager.ErrTokenExpired) {
			ID := refreshTokenClaims.UserID

			auth, err := a.authPostgresRepo.GetUserAuth(ctx, ID)
			if err != nil {
				return "", ErrAccessDenied
			}

			newAccessToken, err := a.authManager.GenerateAccessToken(ctx, ID, auth.Role)
			if err != nil {
				return "", ErrGenerateToken
			}

			return newAccessToken, nil
		}

		return "", ErrAccessDenied
	}

	return accessToken, nil
}

func (a *authenticateManager) SendResetPasswordVerification(ctx context.Context, email string) (token string, _ error) {
	user, err := a.userPostgresRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	auth, err := a.authPostgresRepo.GetUserAuth(ctx, user.ID)
	if err != nil {
		return "", err
	}

	if !auth.EmailVerified {
		return "", ErrEmailNotVerified
	}

	if auth.FailedLoginAttempts >= MaximumFailedLoginAttempts {
		return "", fmt.Errorf("%w until: %v", ErrAccountLocked, auth.AccountLockedUntil)
	}

	resetPasswordToken, err := a.authManager.GenerateResetPasswordToken(ctx, user.ID)
	if err != nil {
		return "", ErrGenerateToken
	}

	err = a.emailService.SendResetPasswordEmail(email, "example.com", user.FirstName, "10")
	if err != nil {
		return "", err
	}

	return resetPasswordToken, nil
}

func (a *authenticateManager) SubmitResetPassword(ctx context.Context, token string, newPassword string) error {
	tokenClaims, err := a.authManager.DecodeResetPasswordToken(ctx, token)
	if err != nil {
		return ErrAccessDenied
	}

	ID := convertString2Uint(tokenClaims.UserID)

	auth, err := a.authPostgresRepo.GetUserAuth(ctx, ID)
	if err != nil {
		return ErrAccessDenied
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return ErrHashingPassword
	}

	err = a.authPostgresRepo.ChangePassword(ctx, auth.ID, newPasswordHash)
	if err != nil {
		return ErrChangePassword
	}

	return nil
}

func (a *authenticateManager) DeleteAccount(ctx context.Context,userID uint, password string) error {	
	auth, err := a.authPostgresRepo.GetUserAuth(ctx,userID)
	if err != nil {
		return ErrNotFound
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.HashedPassword)
	if !validPassword {
		return ErrDeleteUser
	}

	err = a.authPostgresRepo.DeleteByID(ctx, userID)
	if err != nil {
		return ErrDeleteUser
	}

	err = a.userPostgresRepo.DeleteByID(ctx, userID)
	if err != nil {
		return ErrDeleteUser
	}

	return nil
}

func (a *authenticateManager) DestroyRefreshToken(ctx context.Context, refreshToken string) error {
	err := a.authManager.DestroyRefreshToken(ctx, refreshToken)
	if err != nil {
		return ErrNotFound
	}
	return nil
}

//you can use it for logout 
func (a *authenticateManager) SetAccessTokenIntoBlacklist(ctx context.Context, accessToken string,accessTokenExpr time.Duration) error {
	err := a.authManager.SetAccessTokenIntoBlacklist(ctx, accessToken, accessTokenExpr)
	if err != nil {
		return ErrNotFound
	}

	return nil
}

// you can use it in middlewares when someone wants to get in
func (a *authenticateManager) IsAccessTokenBlacklisted(ctx context.Context, accessToken string) bool {
	return a.authManager.IsAccessTokenBlacklisted(ctx, accessToken)
}


func convertString2Uint(strID string) uint {
	ID, err := strconv.Atoi(strID)
	if err != nil {
		panic(err)
	}
	return uint(ID)
}
