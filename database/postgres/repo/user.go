package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"blog/utils"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type userPostgresRepo struct {
	postgresCLI *gorm.DB
}

func NewUserPostgresRepository(postgresCLI *gorm.DB) repository.UserRepository {
	return &userPostgresRepo{
		postgresCLI: postgresCLI,
	}
}
func (u *userPostgresRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	err := u.postgresCLI.WithContext(ctx).Preload("Articles").Preload("Comments").Where("email= ?", email).First(user).Error
	if err != nil {
		return nil, fmt.Errorf("get user by email: %w: %v", repository.ErrDatabase, err)
	}

	return user, nil
}
func (u *userPostgresRepo) GetUserByID(ctx context.Context, ID uint) (*model.User, error) {
	user := &model.User{}
	err := u.postgresCLI.WithContext(ctx).Preload("Articles").Preload("Comments").Preload("Likes.Article.Comments").Preload("Likes.Article.Categories").First(user, ID).Error
	if err != nil {
		return nil, fmt.Errorf("get user by id: user ID:  %d\n%w: %v", ID, repository.ErrDatabase, err)
	}

	return user, nil
}
func (u *userPostgresRepo) Create(ctx context.Context, user *model.User) (*model.User, *gorm.DB, error) {
	tx := NewTx(u.postgresCLI)

	tx = tx.Create(user)
	if tx.Error != nil {

		if utils.CheckErrorForWord(tx.Error, "email") {
			return nil, nil, repository.ErrEmailAlreadyExits
		} else {
			return nil, nil, fmt.Errorf("create user: %v \n%w: %v", user, repository.ErrDatabase, tx.Error)
		}
	}
	//retrun tx for rollback if jwt token can not be set
	return user, tx, nil
}
func (u *userPostgresRepo) UpdateByID(ctx context.Context, ID uint, firstName, lastName, biography string) (*model.User, error) {
	user,err := u.GetUserByID(ctx,ID)
	if err != nil {
		return nil, err
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Biography = biography

	err = u.postgresCLI.WithContext(ctx).Save(user).Error
	if err != nil {
		return nil, fmt.Errorf("update user by id: user ID:%d ,first name:%s ,last name:%s ,biography:%s\n%w: %v",ID,firstName,lastName ,biography,repository.ErrDatabase,err)
	}

	return user, nil
}

func (u *userPostgresRepo) DeleteByID(ctx context.Context, ID uint) error {
	user := new(model.User)

	result := u.postgresCLI.WithContext(ctx).Delete(user, ID)
	if result.Error != nil {
		return fmt.Errorf("delete user by id: user ID:%d \n%w: %v",ID,repository.ErrDatabase,result.Error)
	}

	if result.RowsAffected == 0 {
		return repository.ErrUserNotFound
	}

	return nil
}
