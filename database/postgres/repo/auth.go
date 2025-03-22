package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type authPostgresRepository struct {
	postgresCLI *gorm.DB
}

func NewAuthPostgresRepository(postgresCLI *gorm.DB) repository.AuthRepository {
	return &authPostgresRepository{
		postgresCLI: postgresCLI,
	}
}
func (a *authPostgresRepository) Create(ctx context.Context, authModel *model.Auth) (*model.Auth, error) {
	err := a.postgresCLI.Create(authModel).Error

	if err != nil {
		return nil, fmt.Errorf("create auth: %w: %v", repository.ErrDatabase, err)
	}

	return authModel, nil
}
func (a *authPostgresRepository) GetUserAuth(ctx context.Context, ID uint) (*model.Auth, error) {
	auth := new(model.Auth)

	err := a.postgresCLI.WithContext(ctx).First(auth, ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, repository.ErrAuthNotFound
		}

		return nil, fmt.Errorf("get user auth :auth ID:%d\n%w: %v", ID, repository.ErrDatabase, err)
	}

	return auth, nil
}

func (a *authPostgresRepository) ChangePassword(ctx context.Context, ID uint, hashedPassword string) error {
	authModel, err := a.GetUserAuth(ctx, ID)
	if err != nil {
		return err
	}

	authModel.HashedPassword = hashedPassword
	err = a.postgresCLI.WithContext(ctx).Save(authModel).Error
	if err != nil {
		return fmt.Errorf("change password: auth ID:%d\n%w: %v ", ID, repository.ErrDatabase, err)
	}

	return nil
}
func (a *authPostgresRepository) VerifyEmail(ctx context.Context, ID uint) error {
	auth, err := a.GetUserAuth(ctx, ID)
	if err != nil {
		return err
	}

	auth.EmailVerified = true
	err = a.postgresCLI.WithContext(ctx).Save(auth).Error
	if err != nil {
		return fmt.Errorf("verify email:auth ID:%d\n%w: %v", ID, repository.ErrDatabase, err)
	}

	return nil
}

func (a *authPostgresRepository) IncrementFailedLoginAttempts(ctx context.Context, ID uint) error {
	auth, err := a.GetUserAuth(ctx, ID)
	if err != nil {
		return err
	}

	auth.FailedLoginAttempts += 1
	err = a.postgresCLI.WithContext(ctx).Save(auth).Error
	if err != nil {
		return fmt.Errorf("increment failed login attempts:auth ID:%d \n%w: %v", ID, repository.ErrDatabase, err)
	}

	return nil
}

func (a *authPostgresRepository) ClearFailedLoginAttempts(ctx context.Context, ID uint) error {
	auth, err := a.GetUserAuth(ctx, ID)
	if err != nil {
		return err
	}

	auth.FailedLoginAttempts = 0
	err = a.postgresCLI.WithContext(ctx).Save(auth).Error
	if err != nil {
		return fmt.Errorf("clear failed login attempts: auth ID:%d \n%w: %v", ID, repository.ErrDatabase, err)
	}

	return nil
}

func (a *authPostgresRepository) LockAccount(ctx context.Context, ID uint, lockDuration time.Duration) error {
	auth, err := a.GetUserAuth(ctx, ID)
	if err != nil {
		return err
	}

	now := time.Now()
	now = now.Add(lockDuration)
	auth.AccountLockedUntil = now.Unix()

	err = a.postgresCLI.WithContext(ctx).Save(auth).Error
	if err != nil {
		return fmt.Errorf("lock account: auth ID:%d \n%w: %v", ID, repository.ErrDatabase, err)
	}

	return nil
}
func (a *authPostgresRepository) UnlockAccount(ctx context.Context, ID uint) error {
	auth, err := a.GetUserAuth(ctx, ID)
	if err != nil {
		return err
	}

	auth.AccountLockedUntil = 0

	err = a.postgresCLI.WithContext(ctx).Save(auth).Error
	if err != nil {
		return fmt.Errorf("unlock account:auth ID:%d \n%w: %v", ID, repository.ErrDatabase, err)
	}

	return nil
}

func (a *authPostgresRepository) DeleteByID(ctx context.Context, ID uint) error {
	auth := new(model.Auth)

	result := a.postgresCLI.WithContext(ctx).Delete(auth, ID)
	if result.Error != nil {
		return fmt.Errorf("delete by id:auth ID:%d \n%w: %v", ID, repository.ErrDatabase, result.Error)
	}

	if result.RowsAffected == 0 {
		return repository.ErrAuthNotFound
	}

	return nil
}
