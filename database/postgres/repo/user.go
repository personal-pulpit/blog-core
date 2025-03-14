package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"blog/utils"
	"context"

	"gorm.io/gorm"
)

type userPostgresRepo struct {
	postgresCLI *gorm.DB
}

func NewUserPostgresRepository(postgresCLI *gorm.DB) repository.UserPostgresRepository {
	return &userPostgresRepo{
		postgresCLI: postgresCLI,
	}
}
func (u *userPostgresRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	tx := u.postgresCLI.WithContext(ctx).Preload("Articles").Preload("Comments").Where("email= ?", email).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (u *userPostgresRepo) GetUserByID(ctx context.Context, ID uint) (*model.User, error) {
	user := &model.User{}
	tx := u.postgresCLI.WithContext(ctx).Preload("Articles").Preload("Comments").First(user, ID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (u *userPostgresRepo) Create(ctx context.Context, user *model.User) (*model.User, *gorm.DB, error) {
	tx := NewTx(u.postgresCLI)

	tx = tx.Create(user)
	if tx.Error != nil {
		if utils.CheckErrorForWord(tx.Error, "email") {
			return nil, nil, ErrEmailAlreadyExits
		} else {
			return nil, nil, tx.Error
		}
	}
	//retrun tx for rollback if jwt token can not be set
	return user, tx, nil
}
func (u *userPostgresRepo) UpdateByID(ctx context.Context, ID uint, firstName, lastName, biography string) (*model.User, error) {
	var user = &model.User{}

	err := u.postgresCLI.WithContext(ctx).First(user, ID).Error
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Biography = biography

	err = u.postgresCLI.WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userPostgresRepo) DeleteByID(ctx context.Context, ID uint) error {
	user := new(model.User)

	result := u.postgresCLI.WithContext(ctx).Delete(user, ID)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
