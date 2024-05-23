package repo

import (
	"blog/pkg/data/database"
	"blog/pkg/data/models"
)
func GetUsers()([]models.User,error){
	var users []models.User
	err := database.DB.Find(&users).Error
	return users, err
}
func GetUser(id string) (models.User, error) {
	var u models.User
	err := database.DB.First(&u,id).Error
	return u, err
}

func CreateUser(fristname, lastname, username, password, email, phonenumber string) (models.User, error) {
	var u models.User
	u.Fristname = fristname
	u.Lastname = lastname
	u.Username = username
	u.Password = password
	u.Email = email
	u.PhoneNumber = phonenumber

	err := database.DB.Create(&u).Error

	return u, err
}
func UpdateUserById(id, fristname, lastname, username, password, email, phonenumber string) (models.User,error) {
	u, err := GetUser(id)

	if err != nil {
		return models.User{},err
	}
	u.Fristname = fristname
	u.Lastname = lastname
	u.Username = username
	u.Password = password
	u.Email = email
	u.PhoneNumber = phonenumber

	err = database.DB.Save(&u).Error

	return u,err
}
func DeleteUser(id string) error {
	var u models.User
	err := database.DB.Delete(&u, id).Error

	return err
}
