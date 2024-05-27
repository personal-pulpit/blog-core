package repo

import (
	"blog/pkg/data/database"
	"blog/pkg/data/models"
	"context"
	"fmt"
	"strconv"
)

func GetUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users).Error
	return users, err
}
func GetUser(id string) (models.User, error) {
	var u models.User
	err := database.DB.First(&u, id).Error
	return u, err
}

func CreateUser(firstname, lastname, username, password, email, phonenumber string) (models.User, error) {
	var u models.User
	u.Firstname = firstname
	u.Lastname = lastname
	u.Username = username
	u.Password = password
	u.Email = email
	u.PhoneNumber = phonenumber

	err := database.DB.Create(&u).Error
	if err != nil{
		return models.User{},err
	}
	redisRes := database.Rdb.HMSet(context.Background(), fmt.Sprintf("user:%d", u.Id), map[string]interface{}{
		"firstname":   u.Firstname,
		"lastname":    u.Lastname,
		"username":    u.Username,
		"email":       u.Email,
		"phonenumber": u.PhoneNumber,
		"role":        u.Role,
	})
	return u, redisRes.Err()
}
func UpdateUserById(id, firstname, lastname, username, password, email, phonenumber string) (models.User, error) {
	redisMapRes := database.Rdb.HGetAll(context.Background(), fmt.Sprintf("user:%s", id))
	if redisMapRes.Err() != nil {
		return models.User{}, redisMapRes.Err()
	}
	var u models.User
	for key, value := range redisMapRes.Val() {
		switch key {
		case "firstname":
			u.Firstname = value
		case "lastname":
			u.Lastname = value
		case "username":
			u.Username = value
		case "password":
			u.Password = value
		case "email":
			u.Email = value
		case "phonenumber":
			u.PhoneNumber = value
		}
	}
	err := database.DB.Save(&u).Error
	if err != nil {
		return models.User{}, err
	}
	redisRes := database.Rdb.HMSet(context.Background(), fmt.Sprintf("user:%d", u.Id), map[string]interface{}{
		"firstname":   u.Firstname,
		"lastname":    u.Lastname,
		"username":    u.Username,
		"email":       u.Email,
		"phonenumber": u.PhoneNumber,
	})
	if redisRes.Err() != nil {
		return models.User{}, redisRes.Err()
	}
	return u, err
}
func DeleteUser(id string) error {
	var u models.User
	err := database.DB.Delete(&u, id).Error
	if err != nil {
		return err
	}
	redisRes := database.Rdb.Del(context.Background(), fmt.Sprintf("user:%d", u.Id))
	return redisRes.Err()

}
func VerifyUser(username string) (models.User, error) {
	var u models.User
	err := database.DB.First(&u, "username=?", username).Error
	return u, err
}
func GetUserByIdRedis(id string) (models.User,error) {
	redisMapRes := database.Rdb.HGetAll(context.Background(), fmt.Sprintf("user:%s", id))
	if redisMapRes.Err() != nil {
		return models.User{}, redisMapRes.Err()
	}
	var u models.User
	for key, value := range redisMapRes.Val() {
		switch key {
		case "firstname":
			u.Firstname = value
		case "lastname":
			u.Lastname = value
		case "username":
			u.Username = value
		case "password":
			u.Password = value
		case "email":
			u.Email = value
		case "phonenumber":
			u.PhoneNumber = value
		case "role":
			rolef, _ := strconv.Atoi(value)
			roleu := uint(rolef)
			u.Role = roleu
		}
	}
	return u,nil
}
