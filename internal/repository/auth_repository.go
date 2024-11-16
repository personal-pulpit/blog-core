package repository

import (
	"blog/internal/model"
	"context"
	"time"
)

type AuthPostgresRepository interface {
	Create(ctx context.Context, authModel *model.Auth) (*model.Auth, error)
	GetUserAuth(ctx context.Context,ID uint) (*model.Auth, error)
	ChangePassword(ctx context.Context,ID uint, hashedPassword string) error
	VerifyEmail(ctx context.Context,ID uint) error
	IncrementFailedLoginAttempts(ctx context.Context,ID uint) error
	ClearFailedLoginAttempts(ctx context.Context,ID uint) error
	LockAccount(ctx context.Context,ID uint, lockDuration time.Duration) error
	UnlockAccount(ctx context.Context,ID uint) error
	DeleteByID(ctx context.Context,ID uint) error
}
type AuthRedisRepository interface {
}
