package auth_manager

import (
	"blog/internal/model"
	"blog/utils/random"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/redis/go-redis/v9"
)

var (
	ErrInvalidToken            = errors.New("invalid token")
	ErrInvalidTokenType        = errors.New("invalid token type")
	ErrUnexpectedSigningMethod = errors.New("unexpected token signing method")
	ErrNotFound                = errors.New("not found")
)

var TokenEncodingAlgorithm = jwt.SigningMethodHS512

type TokenType int

const (
	AccessToken TokenType = iota
	RefreshToken
	ResetPassword
	VerifyEmail
)

type AuthManager interface {
	GenerateToken(tokenType TokenType, tokenPayload *TokenClaims, expr time.Duration) (token string, err error)
	DecodeToken(token string, tokenType TokenType) (claims *TokenClaims, err error)
	Destroy(key string) (err error)
	GetOTP(uniqueID string) (otp string, err error)
	SetOTP( uniqueID string, expr time.Duration) (otp string, err error)
}

type AuthManagerOpts struct {
	PrivateKey string
}

// Used as jwt claims
type TokenClaims struct {
	ID        model.ID  `json:"id"`
	Role      model.Role `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

func NewTokenClaims(ID model.ID, role model.Role ,tokenType TokenType) *TokenClaims {
	return &TokenClaims{
		ID:        ID,
		CreatedAt: time.Now(),
		TokenType: tokenType,
	}
}

type authManager struct {
	redisClient *redis.Client
	opts        AuthManagerOpts
}

func NewAuthManager(redisClient *redis.Client, opts AuthManagerOpts) AuthManager {
	return &authManager{redisClient, opts}
}

func (t *authManager) GenerateToken(tokenType TokenType, tokenClaims *TokenClaims, expr time.Duration) (_ string, _ error) {
	token, err := jwt.NewWithClaims(TokenEncodingAlgorithm, tokenClaims).SignedString([]byte(t.opts.PrivateKey))
	if err != nil {
		return "", err
	}

	cmd := t.redisClient.Set(context.TODO(), token, nil, expr)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}

	return token, nil
}

func (t *authManager) DecodeToken(token string, tokenType TokenType) (_ *TokenClaims, _ error) {
	exists,err :=t.redisClient.Exists(context.TODO(),token).Result()
	if err != nil{
		return nil,err
	}
	
	if exists != 1{
		return nil,ErrInvalidToken
	}

	tokenClaims := &TokenClaims{}
	jwtToken, err := jwt.ParseWithClaims(token, tokenClaims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}

			return []byte(t.opts.PrivateKey), nil
		},
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if jwtToken.Valid {
		if tokenClaims.TokenType != tokenType {
			return nil, ErrInvalidTokenType
		}

		return tokenClaims, nil
	}

	return &TokenClaims{}, ErrInvalidToken
}

func (t *authManager) Destroy(key string) (_ error) {
	cmd := t.redisClient.Del(context.TODO(), key)
	if cmd.Err() != nil {
		return cmd.Err()
	}
	return nil
}

func (t *authManager) GetOTP( uniqueID string) (_ string, _ error) {
	result, err := t.redisClient.Get(context.TODO(), uniqueID).Result()
	if err != nil {
		return "", err
	}

	if len(strings.TrimSpace(result)) > 0 {
		return result, nil
	}

	return "", ErrNotFound
}

func (t *authManager) SetOTP(uniqueID string, expr time.Duration) (_ string, _ error) {
	otp := fmt.Sprintf("%d", random.GenerateOTP())

	_, err := t.redisClient.Set(context.TODO(), uniqueID, otp, expr).Result()
	if err != nil {
		return "", err
	}

	return otp, nil
}
