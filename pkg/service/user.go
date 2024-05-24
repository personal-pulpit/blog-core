package service

import (
	"blog/pkg/data/models"
	"blog/pkg/data/repo"
	"blog/utils"
)

func GetUsers() ([]models.User, error) {
	return repo.GetUsers()

}
func GetUser(id string) (models.User, error) {
	return repo.GetUser(id)
}
func CreateUser(fristname, lastname, username, password, email, phonenumber string) (models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}
	return repo.CreateUser(fristname, lastname, username, hashedPassword, email, phonenumber)

}
func VerifyUser(username, password string) (models.User, error) {
	u, err := repo.VerifyUser(username)
	if err != nil {
		return models.User{}, err

	}
	err = utils.CheckPassword(u.Password, password)
	return u, err
}
func UpdateUserById(id, fristname, lastname, username, password, email, phonenumber string) (models.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return models.User{}, err

	}
	return repo.UpdateUserById(id, fristname, lastname, username, hashedPassword, email, phonenumber)

}
func DeleteUser(id string) error {
	return repo.DeleteUser(id)
}
