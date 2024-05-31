package repo

import (
	"blog/pkg/data/database"
	"blog/pkg/data/models"
	"blog/utils"
	"errors"

	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB  *gorm.DB
	RDB *redis.Client
}

var (
	ErrEmailAlreadyExits       = errors.New("email already exits")
	ErrUsernameAlreadyExits    = errors.New("username already exits")
	ErrPhoneNumberAlreadyExits = errors.New("phone number already exits")
	ErrUserNotFound            = errors.New("user not found")
	ErrUsernameOrPasswordWrong      = errors.New("username or password wrong")
)

func NewUserRepo() *UserRepo {
	return &UserRepo{
		DB:  database.DB,
		RDB: database.Rdb,
	}
}
func (ur *UserRepo) CreateUser(firstname, lastname, biography, username, password, email, phonenumber string) (models.User, error) {
	var u models.User
	u.Firstname = firstname
	u.Lastname = lastname
	u.Biography = biography
	u.Username = username
	u.Email = email
	u.PhoneNumber = phonenumber
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return models.User{}, err
	}
	u.Password = hashedPassword
	err = ur.DB.Create(&u).Error
	if err != nil {
		if utils.CheckErrorForWord(err, "email") {
			return models.User{}, ErrEmailAlreadyExits
		} else if utils.CheckErrorForWord(err, "username") {
			return models.User{}, ErrUsernameAlreadyExits
		} else if utils.CheckErrorForWord(err, "phone_number") {
			return models.User{}, ErrPhoneNumberAlreadyExits
		} else {
			return models.User{}, err
		}
	}
	return ur.CreateChache(u)
}
func (ur *UserRepo) UpdateUserById(id, firstname, lastname, biography, username string) (models.User, error) {
	var u models.User
	err := ur.DB.First(&u, id).Error
	if err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return u,ErrUserNotFound
		}
		return u, err
	}
	u.Firstname = firstname
	u.Lastname = lastname
	u.Biography = biography
	u.Username = username
	err = ur.DB.Save(&u).Error
	if err != nil {
		return models.User{}, err
	}
	err = ur.DeleteChacheByIdRedis(id)
	if err != nil {
		return models.User{}, err
	}
	u, err = ur.CreateChache(u)
	if err != nil {
		return models.User{}, err
	}
	return u, err
}
func (ur *UserRepo) DeleteUser(id string) error {
	var u models.User
	err := ur.DB.Delete(&u, id).Error
	if err != nil {
		return err
	}
	id = strconv.Itoa(int(u.Id))
	err = ur.DeleteChacheByIdRedis(id)
	return err
}
func (ur *UserRepo) VerifyUser(username, password string) (models.User, error) {
	var u models.User
	err := ur.DB.First(&u, "username=?", username).Error
	if err != nil {
		if errors.Is(err,gorm.ErrRecordNotFound){
			return u,ErrUserNotFound
		}
		return u, err
	}
	err = utils.CheckPassword(password, u.Password)
	if err != nil {
		if utils.CheckErrorForWord(err,"crypto/bcrypt"){
			return u,ErrUsernameOrPasswordWrong
		}
		return u, err
	}
	u, err = ur.CreateChache(u)
	return u, err
}
func (ur *UserRepo) GetUserByIdRedis(id string) (map[string]string, error) {
	exists := ur.RDB.Exists(context.Background(), fmt.Sprintf("user:%s", id))
	if exists.Val() == 0 {
		return map[string]string{}, ErrUserNotFound
	}
	redisMapRes := ur.RDB.HGetAll(context.Background(), fmt.Sprintf("user:%s", id))
	if redisMapRes.Err() != nil {
		return map[string]string{}, redisMapRes.Err()
	}
	return redisMapRes.Val(), nil
}
func (ur *UserRepo) CreateChache(u models.User) (models.User, error) {
	redisRes := database.Rdb.HMSet(context.Background(), fmt.Sprintf("user:%d", u.Id), map[string]interface{}{
		"firstname":   u.Firstname,
		"lastname":    u.Lastname,
		"biography":   u.Biography,
		"username":    u.Username,
		"email":       u.Email,
		"phonenumber": u.PhoneNumber,
		"role":        u.Role,
		"createdAt":   u.CreatedAt,
		"updatedAt":   u.UpdatedAt,
	})
	return u, redisRes.Err()
}
func (ur *UserRepo) DeleteChacheByIdRedis(id string) error {
	redisRes := database.Rdb.Del(context.Background(), fmt.Sprintf("user:%s", id))
	return redisRes.Err()
}
func (ur *UserRepo) GetUsersRedis() ([]map[string]string, error) {
	var users []map[string]string
	keys, err := ur.RDB.Keys(context.Background(), "user:*").Result()
	if err != nil {
		return users, err
	}
	for _, key := range keys {
		userMap, err := ur.RDB.HGetAll(context.Background(), key).Result()
		if err != nil {
			return []map[string]string{}, err
		}
		users = append(users, userMap)
	}
	return users, nil
}
