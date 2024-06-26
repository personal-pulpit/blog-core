package repository

import (
	"blog/internal/model"
	"time"
)

type AuthPostgresRepository interface {
	Create(authModel *model.Auth) (*model.Auth, error)
	GetUserAuth(ID model.ID) (*model.Auth, error)
	ChangePassword(ID model.ID, hashedPassword string) error
	VerifyEmail(ID model.ID) error
	IncrementFailedLoginAttempts(ID model.ID) error
	ClearFailedLoginAttempts(ID model.ID) error
	LockAccount(ID model.ID, lockDuration time.Duration) error
	UnlockAccount(ID model.ID) error
	DeleteByID(ID model.ID) error
}
type AuthRedisRepository interface {
}
