package mysql_repository

import (
	database "blog/database/mysql"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/utils"
	"errors"

	"gorm.io/gorm"
)

type userMysqlRepo struct {
	mysqlClient *gorm.DB
}

func NewUserMysqlRepository() repository.UserMysqlRepository {
	return &userMysqlRepo{
		mysqlClient: database.GetMysqlDB(),
	}
}
func (u *userMysqlRepo) Create(firstname, lastname, biography, username, password, email, phonenumber string) (model.User, *gorm.DB, error) {
	var user model.User
	user.Firstname = firstname
	user.Lastname = lastname
	user.Biography = biography
	user.Username = username
	user.Email = email
	user.PhoneNumber = phonenumber
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return user, nil, err
	}
	user.Password = hashedPassword
	tx := NewTx(u.mysqlClient)
	err = tx.Create(&u).Error
	if err != nil {
		if utils.CheckErrorForWord(err, "email") {
			return user, nil, ErrEmailAlreadyExits
		} else if utils.CheckErrorForWord(err, "username") {
			return user, nil, ErrUsernameAlreadyExits
		} else if utils.CheckErrorForWord(err, "phone_number") {
			return user, nil, ErrPhoneNumberAlreadyExits
		} else {
			return user, nil, err
		}
	}
	txj := tx.Commit()
	//retrun tx for rollback if jwt token can not be set
	return user, txj, nil
}
func (u *userMysqlRepo) UpdateByID(ID, firstname, lastname, biography, username string) (model.User, error) {
	var user model.User
	tx := NewTx(u.mysqlClient)
	err := tx.First(&u, ID).Error
	if err != nil {
		return user, err
	}
	user.Firstname = firstname
	user.Lastname = lastname
	user.Biography = biography
	user.Username = username
	err = tx.Save(&u).Error
	if err != nil {
		tx.Rollback()
		return user, err
	}
	return user, nil
}
func (u *userMysqlRepo) DeleteByID(ID string) error {
	tx := NewTx(u.mysqlClient)
	err := tx.Delete(&u, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
	}
	return nil
}
func (u *userMysqlRepo) Verify(username, password string) (model.User, error) {
	var user model.User
	err := u.mysqlClient.First(&u, "username=?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, ErrUserNotFound
		}
		return user, err
	}
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		if utils.CheckErrorForWord(err, "crypto/bcrypt") {
			return user, ErrUsernameOrPasswordWrong
		}
		return user, err
	}
	return user, nil
}
