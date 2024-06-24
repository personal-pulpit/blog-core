package mysql_repository

import (
	"blog/internal/model"
	"blog/internal/repository"
	"blog/utils"
	"errors"

	"gorm.io/gorm"
)

type userMysqlRepo struct {
	mysqlClient *gorm.DB
}

func NewUserMysqlRepository(mysqlCLI *gorm.DB) repository.UserMysqlRepository {
	return &userMysqlRepo{
		mysqlClient: mysqlCLI,
	}
}
func (u *userMysqlRepo) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	tx := u.mysqlClient.Where("email= ?", email).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (u *userMysqlRepo) GetUserByID(ID model.ID) (*model.User, error) {
	user := &model.User{}
	tx := u.mysqlClient.Where("id= ?", ID).First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (u *userMysqlRepo) Create(user *model.User) (*model.User, *gorm.DB, error) {
	tx := u.mysqlClient.Create(&u)
	if tx.Error != nil {
		if utils.CheckErrorForWord(tx.Error, "email") {
			return user, nil, ErrEmailAlreadyExits
		} else if utils.CheckErrorForWord(tx.Error, "username") {
			return user, nil, ErrUsernameAlreadyExits
		} else if utils.CheckErrorForWord(tx.Error, "phone_number") {
			return user, nil, ErrPhoneNumberAlreadyExits
		} else {
			return user, nil, tx.Error
		}
	}
	//retrun tx for rollback if jwt token can not be set
	return user, tx, nil
}
func (u *userMysqlRepo) UpdateByID(ID, FirstName, lastname, biography string) (*model.User, error) {
	var user = &model.User{}
	err := u.mysqlClient.First(&u, ID).Error
	if err != nil {
		return user, err
	}
	user.FirstName = FirstName
	user.LastName = lastname
	user.Biography = biography
	err = u.mysqlClient.Save(&u).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
func (u *userMysqlRepo) DeleteByID(ID model.ID) error {
	tx := NewTx(u.mysqlClient)
	err := tx.Delete(&u, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
	}
	return nil
}
