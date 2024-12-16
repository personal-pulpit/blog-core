package auth_manager

import (
	"blog/config"
	"blog/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	auth_manager "github.com/tahadostifam/go-auth-manager"
)

type AuthManager interface {
	GenerateAccessToken(ctx context.Context, userID uint, role model.Role) (accessToken string, _ error)
	DecodeAccessToken(ctx context.Context, accessToken string) (*AccessTokenClaims, error)
	GenerateRefreshToken(ctx context.Context, userID uint, ipAddress, userAgent string) (refreshToken string, _ error)
	DecodeRefreshToken(ctx context.Context, RefreshToken string) (*RefreshTokenClaims, error)
	GenerateResetPasswordToken(ctx context.Context, userID uint) (resetPasswordToken string, _ error)
	DecodeResetPasswordToken(ctx context.Context, resetPasswordToken string) (*ResetPasswordTokenClaims, error)
	GenerateVerificationCode(ctx context.Context, key string) (verificationCode string, _ error)
	CompareVerificationCode(ctx context.Context, key string, verificationCode string) (bool, error)
	DestroyPailToken(ctx context.Context, key string) (_ error)
	DestroyRefreshToken(ctx context.Context, key string) error
}
type AccessTokenClaims struct {
	UserID   string
	Role     string
	CreateAt time.Time
}

type RefreshTokenClaims struct {
	IPAddress  string
	UserAgent  string
	UserID     uint
	LoggedInAt time.Duration
}

type ResetPasswordTokenClaims struct {
	UserID   string
	CreateAt time.Time
}

const (
	ResetPasswordTokenExpr = time.Second * 360   // 5 minutes
	AccessTokenExpr        = time.Minute * 45    // 45 minutes
	RefreshTokenExpr       = time.Hour * 24 * 14 // 2 weeks
	VerificationCodeExpr   = time.Minute * 2     // 2 minutes
	VerificationCodeLength = 6
)

type authManger struct {
	authManger auth_manager.AuthManager
}

func NewAuthManger(redisClient *redis.Client, jwtConfigs config.Jwt) AuthManager {
	mainAuthManger := auth_manager.NewAuthManager(redisClient, auth_manager.AuthManagerOpts{
		PrivateKey: jwtConfigs.Secret,
	})

	authManger := &authManger{
		authManger: mainAuthManger,
	}

	return authManger

}

func (a *authManger) GenerateAccessToken(ctx context.Context, userID uint, role model.Role) (accessToken string, _ error) {
	accessToken, err := a.authManger.GenerateAccessToken(ctx, fmt.Sprint(userID), fmt.Sprint(role), AccessTokenExpr)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (a *authManger) DecodeAccessToken(ctx context.Context, accessToken string) (*AccessTokenClaims, error) {
	accessTokenPayload, err := a.authManger.DecodeAccessToken(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	accessTokenClaims := &AccessTokenClaims{
		UserID: accessTokenPayload.Payload.UUID,
		Role:   accessTokenPayload.Payload.Role,
	}

	return accessTokenClaims, nil
}

func (a *authManger) GenerateRefreshToken(ctx context.Context, userID uint, ipAddress, userAgent string) (refreshToken string, _ error) {
	refreshTokenClaims, err := a.authManger.GenerateRefreshToken(ctx, &auth_manager.RefreshTokenPayload{
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		LoggedInAt: time.Duration(time.Now().UnixMilli()),
		UserID:     userID,
	},
		RefreshTokenExpr)
	if err != nil {
		return "", err
	}

	return refreshTokenClaims, nil
}

func (a *authManger) DecodeRefreshToken(ctx context.Context, RefreshToken string) (*RefreshTokenClaims, error) {
	refreshTokenPayload, err := a.authManger.DecodeRefreshToken(ctx, RefreshToken)
	if err != nil {
		return nil, err
	}

	refreshTokenClaims := &RefreshTokenClaims{
		IPAddress:  refreshTokenPayload.IPAddress,
		UserAgent:  refreshTokenPayload.UserAgent,
		UserID:     refreshTokenPayload.UserID,
		LoggedInAt: refreshTokenPayload.LoggedInAt,
	}

	return refreshTokenClaims, err
}

func (a *authManger) GenerateResetPasswordToken(ctx context.Context, userID uint) (resetPasswordToken string, _ error) {
	resetPasswordToken, err := a.authManger.GeneratePlainToken(ctx, auth_manager.ResetPassword, &auth_manager.TokenPayload{UUID: fmt.Sprintf("%d", userID)}, ResetPasswordTokenExpr)
	if err != nil {
		return "", err
	}

	return resetPasswordToken, nil
}

func (a *authManger) DecodeResetPasswordToken(ctx context.Context, resetPasswordToken string) (*ResetPasswordTokenClaims, error) {
	resetPasswordTokePayload, err := a.authManger.DecodePlainToken(ctx, resetPasswordToken, auth_manager.ResetPassword)
	if err != nil {
		return nil, err
	}
	resetPasswordTokeClaims := &ResetPasswordTokenClaims{
		UserID:   resetPasswordTokePayload.UUID,
		CreateAt: resetPasswordTokePayload.CreatedAt,
	}

	return resetPasswordTokeClaims, nil
}

func (a *authManger) GenerateVerificationCode(ctx context.Context, key string) (verificationCode string, _ error) {
	verificationCode, err := a.authManger.GenerateVerificationCode(ctx, key, VerificationCodeLength, VerificationCodeExpr)
	if err != nil {
		return "", err
	}

	return verificationCode, nil
}

func (a *authManger) CompareVerificationCode(ctx context.Context, key string, verificationCode string) (bool, error) {
	isEqual, err := a.authManger.CompareVerificationCode(ctx, key, verificationCode)
	if err != nil {
		return false, err
	}

	return isEqual, nil
}

func (a *authManger) DestroyPailToken(ctx context.Context, key string) (_ error) {
	err := a.authManger.DestroyPlainToken(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (a *authManger) DestroyRefreshToken(ctx context.Context, key string) error {
	err := a.authManger.TerminateRefreshTokens(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
