package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
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
func (a *authPostgresRepository) Create(authModel *model.Auth) (*model.Auth, error) {
	tx := a.postgresCLI.Create(authModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return authModel, nil
}
func (a *authPostgresRepository) GetUserAuth(ID model.ID) (*model.Auth, error) {
	auth := &model.Auth{}
	tx := a.postgresCLI.Where("id= ?", ID).First(auth)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return auth, nil
}
func (a *authPostgresRepository) ChangePassword(ID model.ID, hashedPassword string) error {
	authModel, err := a.GetUserAuth(ID)
	if err != nil {
		return err
	}
	authModel.HashedPassword = hashedPassword
	tx := a.postgresCLI.Save(authModel)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (a *authPostgresRepository) VerifyEmail(ID model.ID) error {
	auth, err := a.GetUserAuth(ID)
	if err != nil {
		return err
	}
	auth.EmailVerified = true
	tx := a.postgresCLI.Save(auth)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *authPostgresRepository) IncrementFailedLoginAttempts(ID model.ID) error {
	auth, err := a.GetUserAuth(ID)
	if err != nil {
		return err
	}
	auth.FailedLoginAttempts += 1
	tx := a.postgresCLI.Save(auth)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *authPostgresRepository) ClearFailedLoginAttempts(ID model.ID) error {
	auth, err := a.GetUserAuth(ID)
	if err != nil {
		return err
	}
	auth.FailedLoginAttempts = 0
	tx := a.postgresCLI.Save(auth)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *authPostgresRepository) LockAccount(ID model.ID, lockDuration time.Duration) error {
	auth, err := a.GetUserAuth(ID)
	if err != nil {
		return err
	}
	now := time.Now()
	now = now.Add(lockDuration)
	auth.AccountLockedUntil = now.Unix()
	tx := a.postgresCLI.Save(auth)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (a *authPostgresRepository) UnlockAccount(ID model.ID) error {
	auth, err := a.GetUserAuth(ID)
	if err != nil {
		return err
	}
	auth.AccountLockedUntil = 0
	tx := a.postgresCLI.Save(auth)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (a *authPostgresRepository) DeleteByID(ID model.ID) error {
	authModel, err := a.GetUserAuth(ID)
	if err != nil {
		return err
	}
	tx := a.postgresCLI.Delete(&authModel)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
