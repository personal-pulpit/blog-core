package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"time"

	"gorm.io/gorm"
)

type authPostgresRepository struct {
	postgresCLI *gorm.DB
}

func NewAuthPostgresRepository(postgresCLI *gorm.DB) repository.AuthPostgresRepository {
	return &authPostgresRepository{
		postgresCLI: postgresCLI,
	}
}
func (a *authPostgresRepository) Create(ctx context.Context,authModel *model.Auth) (*model.Auth, error) {
	tx := a.postgresCLI.Create(authModel)

	if tx.Error != nil {
		return nil, tx.Error
	}
	
	return authModel, nil
}
func (a *authPostgresRepository) GetUserAuth(ctx context.Context,ID uint) (*model.Auth, error) {
	auth := new(model.Auth)

	tx := a.postgresCLI.WithContext(ctx).First(auth,ID)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return auth, nil
}
func (a *authPostgresRepository) ChangePassword(ctx context.Context,ID uint, hashedPassword string) error {
	authModel, err := a.GetUserAuth(ctx,ID)
	if err != nil {
		return err
	}

	authModel.HashedPassword = hashedPassword
	tx := a.postgresCLI.WithContext(ctx).Save(authModel)
	if tx.Error != nil {
		return tx.Error
	}
	
	return nil
}
func (a *authPostgresRepository) VerifyEmail(ctx context.Context,ID uint) error {
	auth, err := a.GetUserAuth(ctx,ID)
	if err != nil {
		return err
	}

	auth.EmailVerified = true
	tx := a.postgresCLI.WithContext(ctx).Save(auth)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (a *authPostgresRepository) IncrementFailedLoginAttempts(ctx context.Context,ID uint) error {
	auth, err := a.GetUserAuth(ctx,ID)
	if err != nil {
		return err
	}

	auth.FailedLoginAttempts += 1
	tx := a.postgresCLI.WithContext(ctx).Save(auth)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (a *authPostgresRepository) ClearFailedLoginAttempts(ctx context.Context,ID uint) error {
	auth, err := a.GetUserAuth(ctx,ID)
	if err != nil {
		return err
	}

	auth.FailedLoginAttempts = 0
	tx := a.postgresCLI.WithContext(ctx).Save(auth)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (a *authPostgresRepository) LockAccount(ctx context.Context,ID uint, lockDuration time.Duration) error {
	auth, err := a.GetUserAuth(ctx,ID)
	if err != nil {
		return err
	}

	now := time.Now()
	now = now.Add(lockDuration)
	auth.AccountLockedUntil = now.Unix()

	tx := a.postgresCLI.WithContext(ctx).Save(auth)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
func (a *authPostgresRepository) UnlockAccount(ctx context.Context,ID uint) error {
	auth, err := a.GetUserAuth(ctx,ID)
	if err != nil {
		return err
	}

	auth.AccountLockedUntil = 0
	tx := a.postgresCLI.WithContext(ctx).Save(auth)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (a *authPostgresRepository) DeleteByID(ctx context.Context,ID uint) error {
	auth := new(model.Auth)

	result := a.postgresCLI.WithContext(ctx).Delete(auth, ID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
