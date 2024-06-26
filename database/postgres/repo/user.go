package postgres_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"blog/utils"
	"errors"

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
func (u *userPostgresRepo) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	tx := u.postgresCLI.Where("email= ?", email).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (u *userPostgresRepo) GetUserByID(ID model.ID) (*model.User, error) {
	user := &model.User{}
	tx := u.postgresCLI.Where("id= ?", ID).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (u *userPostgresRepo) Create(user *model.User) (*model.User, *gorm.DB, error) {
	tx := u.postgresCLI.Create(&user)
	if tx.Error != nil {
		if utils.CheckErrorForWord(tx.Error, "email") {
			return user, nil, ErrEmailAlreadyExits
		} else {
			return user, nil, tx.Error
		}
	}
	//retrun tx for rollback if jwt token can not be set
	return user, tx, nil
}
func (u *userPostgresRepo) UpdateByID(ID, FirstName, lastname, biography string) (*model.User, error) {
	var user = &model.User{}
	err := u.postgresCLI.First(&u, ID).Error
	if err != nil {
		return user, err
	}
	user.FirstName = FirstName
	user.LastName = lastname
	user.Biography = biography
	err = u.postgresCLI.Save(&u).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func (u *userPostgresRepo) DeleteByID(ID model.ID) error {
	tx := NewTx(u.postgresCLI)
	err := tx.Delete(&u, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
	}
	return nil
}
