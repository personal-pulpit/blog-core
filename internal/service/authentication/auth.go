package authentication

import (
	"blog/api/helpers"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/pkg/auth_manager"
	email "blog/pkg/email_manager"
	"blog/utils/hash"
	"blog/utils/random"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
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

var mutex = new(sync.Mutex)

type AuthService interface {
	Register(ctx context.Context, firstName, lastName, email, biography, password string) (*model.User, error)
	Login(ctx context.Context, email string, password string) (*model.User, string, string, error)
	VerifyEmail(ctx context.Context, otp string, userID uint) error
	SendResetPasswordVerification(ctx context.Context, email string) (string, error)
	SubmitResetPassword(ctx context.Context, token string, newPassword string) error
	ChangePassword(ctx context.Context, userID uint, oldPassword string, newPassword string) error
	Authenticate(ctx context.Context, accessToken string) (*model.User, error)
	RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error)
	SetAccessTokenIntoBlacklist(ctx context.Context, accessToken string, accessTokenExpr time.Duration) error
	IsAccessTokenBlacklisted(ctx context.Context, accessToken string) bool
	DestroyRefreshToken(ctx context.Context, token string) error
}
type authenticateManager struct {
	uniqueID         string
	userPostgresRepo repository.UserRepository
	authPostgresRepo repository.AuthRepository
	authManager      auth_manager.AuthManager
	hashManager      *hash.HashManager
	emailService     email.EmailService
}

func NewAuthenticateService(authRepo repository.AuthRepository, userRepo repository.UserRepository, authManager auth_manager.AuthManager, hashManager *hash.HashManager, emailService email.EmailService) AuthService {
	return &authenticateManager{
		authPostgresRepo: authRepo,
		userPostgresRepo: userRepo,
		authManager:      authManager,
		hashManager:      hashManager,
		emailService:     emailService,
	}
}

func (a *authenticateManager) Register(ctx context.Context, firstName, lastName, email, biography, password string) (*model.User, error) {
	userModel := model.NewUser(firstName, lastName, email, biography)

	savedUser, tx, err := a.userPostgresRepo.Create(ctx, userModel)
	if err != nil {
		return nil, err
	}

	passwordHash, err := a.hashManager.HashPassword(password)

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("Register: %w: %v", ErrHashPassword, err)
	}

	authModel := model.NewAuth(savedUser.ID, passwordHash, model.UserRole)

	_, err = a.authPostgresRepo.Create(ctx, authModel)

	if err != nil {
		tx.Rollback()

		return nil, err
	}

	mutex.Lock()
	uniqueID, err := random.GenerateUniqueID()
	if err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("Register: %w: %v", ErrGenerateUniqueID, err)
	}

	a.uniqueID = fmt.Sprintf("%d", uniqueID)
	mutex.Unlock()

	emailVerificationCode, err := a.authManager.GenerateVerificationCode(ctx, a.uniqueID)

	if err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("Register: %w: %v", ErrGenerateVerificationCode, err)
	}

	err = a.emailService.SendVerificationEmail(userModel.Email, emailVerificationCode)

	if err != nil {
		tx.Rollback()

		return nil, fmt.Errorf("Register: %w: %v", ErrSendVerificationCode, err)

	}

	tx.Commit()

	return savedUser, nil
}

func (a *authenticateManager) VerifyEmail(ctx context.Context, verificationCode string, userID uint) error {
	isValid, err := a.authManager.CompareVerificationCode(ctx, a.uniqueID, verificationCode)

	if err != nil {
		return fmt.Errorf("%w: %v", ErrVerifyEmail, err)
	}

	if !isValid {
		return fmt.Errorf("%w: %w", ErrVerifyEmail, ErrCodeIsInvalid)
	}

	err = a.authPostgresRepo.VerifyEmail(ctx, userID)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrVerifyEmail, err)
	}

	return nil
}
func (a *authenticateManager) Login(ctx context.Context, email string, password string) (*model.User, string, string, error) {
	userModel, err := a.userPostgresRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}

	auth, err := a.authPostgresRepo.GetUserAuth(ctx, userModel.ID)
	if err != nil {
		return nil, "", "", err
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
				return nil, "", "", err
			}

			err = a.authPostgresRepo.ClearFailedLoginAttempts(ctx, auth.ID)
			if err != nil {
				return nil, "", "", err
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
			return nil, "", "", fmt.Errorf("%w : %w", ErrLockAccount, err)
		}
	}

	validPassword := a.hashManager.CheckPasswordHash(password, auth.HashedPassword)
	if !validPassword {
		err = a.authPostgresRepo.IncrementFailedLoginAttempts(ctx, userModel.ID)
		if err != nil {
			return nil, "", "", err
		}

		return nil, "", "", ErrInvalidEmailOrPassword
	}

	accessToken, err := a.authManager.GenerateAccessToken(ctx, userModel.ID, auth.Role)
	if err != nil {
		return nil, "", "", fmt.Errorf("%w : %w", ErrGenerateToken, err)
	}

	refreshToken, err := a.authManager.GenerateRefreshToken(ctx, userModel.ID, "not implemented", "not implemented")
	if err != nil {
		return nil, "", "", fmt.Errorf("%w : %w", ErrGenerateToken, err)
	}

	err = a.authPostgresRepo.ClearFailedLoginAttempts(ctx, auth.ID)
	if err != nil {
		return nil, "", "", err
	}

	err = a.emailService.SendWelcomeEmail(email, userModel.FirstName)
	if err != nil {
		return nil, "", "", fmt.Errorf("%w : %w", ErrSendWelcome, err)
	}

	return userModel, accessToken, refreshToken, nil
}
func (a *authenticateManager) Authenticate(ctx context.Context, accessToken string) (*model.User, error) {

	tokenClaims, err := a.authManager.DecodeAccessToken(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrDecodeAccessToken, err)
	}

	if len(strings.TrimSpace(string(tokenClaims.UserID))) == 0 || len(strings.TrimSpace(tokenClaims.Role)) == 0 {
		return nil, ErrAccessDenied
	} else if !(tokenClaims.Role == fmt.Sprintf("%d", model.UserRole)) && !(tokenClaims.Role == fmt.Sprintf("%d", model.AdminRole)) {
		return nil, ErrAccessDenied
	}

	ID, err := helpers.StringToInt(tokenClaims.UserID)
	if err != nil {
		return nil, fmt.Errorf("%s : %w", "Invalid ID format", err)
	}

	user, err := a.userPostgresRepo.GetUserByID(ctx, uint(ID))
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (a *authenticateManager) ChangePassword(ctx context.Context, userID uint, oldPassword string, newPassword string) error {
	auth, err := a.authPostgresRepo.GetUserAuth(ctx, userID)
	if err != nil {
		return err
	}

	validPassword := a.hashManager.CheckPasswordHash(oldPassword, auth.HashedPassword)
	if !validPassword {
		return fmt.Errorf("%w : %w", ErrCheckPasswordHash, err)
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("%w : %w", ErrHashPassword, err)
	}

	err = a.authPostgresRepo.ChangePassword(ctx, userID, newPasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (a *authenticateManager) RefreshToken(ctx context.Context, refreshToken string, accessToken string) (string, error) {
	refreshTokenClaims, err := a.authManager.DecodeRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", fmt.Errorf("%w : %w", ErrDecodeRefreshToken, err)
	}

	_, err = a.authManager.DecodeAccessToken(ctx, accessToken)

	if err != nil {
		if errors.Is(err, auth_manager.ErrTokenExpired) {
			ID := refreshTokenClaims.UserID

			auth, err := a.authPostgresRepo.GetUserAuth(ctx, ID)
			if err != nil {
				return "", err
			}

			newAccessToken, err := a.authManager.GenerateAccessToken(ctx, ID, auth.Role)
			if err != nil {
				return "", fmt.Errorf("%w : %w", ErrGenerateToken, err)
			}

			return newAccessToken, nil
		}

		return "", fmt.Errorf("%w : %w", ErrRefreshToken, err)
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
		return "", fmt.Errorf("%w : %w", ErrGenerateToken, err)
	}

	err = a.emailService.SendResetPasswordEmail(email, "example.com", user.FirstName, "10")
	if err != nil {
		return "", fmt.Errorf("%w : %w", ErrSendResetPassword, err)
	}

	return resetPasswordToken, nil
}

func (a *authenticateManager) SubmitResetPassword(ctx context.Context, token string, newPassword string) error {
	tokenClaims, err := a.authManager.DecodeResetPasswordToken(ctx, token)
	if err != nil {
		return err
	}

	ID, err := helpers.StringToInt(tokenClaims.UserID)
	if err != nil {
		return fmt.Errorf("%s : %w", "Invalid ID format", err)
	}

	auth, err := a.authPostgresRepo.GetUserAuth(ctx, uint(ID))
	if err != nil {
		return err
	}

	newPasswordHash, err := a.hashManager.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("%w : %w", ErrHashPassword, err)
	}

	err = a.authPostgresRepo.ChangePassword(ctx, auth.ID, newPasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (a *authenticateManager) DestroyRefreshToken(ctx context.Context, refreshToken string) error {
	err := a.authManager.DestroyRefreshToken(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("%w : %w", ErrDestroyRefreshToken, err)
	}
	return nil
}

// you can use it for logout
func (a *authenticateManager) SetAccessTokenIntoBlacklist(ctx context.Context, accessToken string, accessTokenExpr time.Duration) error {
	err := a.authManager.SetAccessTokenIntoBlacklist(ctx, accessToken, accessTokenExpr)
	if err != nil {
		return fmt.Errorf("%w : %w", ErrSetTokenIntoBlackList, err)
	}

	return nil
}

// you can use it in middlewares when someone wants to get in
func (a *authenticateManager) IsAccessTokenBlacklisted(ctx context.Context, accessToken string) bool {
	return a.authManager.IsAccessTokenBlacklisted(ctx, accessToken)
}
