package db

import (
	"blog/pkg/data/database"
	"blog/pkg/data/models"
	"blog/pkg/logging"
	"blog/utils"
	"errors"

	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB     *gorm.DB
	RDB    *redis.Client
	Logger logging.ZapLogger
}

var (
	ErrEmailAlreadyExits       = errors.New("email already exits")
	ErrUsernameAlreadyExits    = errors.New("username already exits")
	ErrPhoneNumberAlreadyExits = errors.New("phone number already exits")
	ErrUserNotFound            = errors.New("user not found")
	ErrUsernameOrPasswordWrong = errors.New("username or password wrong")
)

func NewUserRepo() *UserRepo {
	return &UserRepo{
		DB:     database.DB,
		RDB:    database.Rdb,
		Logger: logging.MyLogger,
	}
}
func (ur *UserRepo) GetAll() ([]map[string]string, error) {
	var users []map[string]string
	keys, err := ur.RDB.Keys(context.Background(), "user:*").Result()
	if err != nil {
		ur.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
		return users, err
	}
	for _, key := range keys {
		userMap, err := ur.RDB.HGetAll(context.Background(), key).Result()
		if err != nil {
			ur.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
			return []map[string]string{}, err
		}
		users = append(users, userMap)
	}
	ur.Logger.Info(logging.Redis, logging.Get, "", nil)
	return users, nil
}
func (ur *UserRepo) GetById(id string) (map[string]string, error) {
	exists := ur.RDB.Exists(context.Background(), fmt.Sprintf("user:%s", id))
	if exists.Val() == 0 {
		ur.Logger.Error(logging.Redis, logging.Get, ErrUserNotFound.Error(), nil)
		return map[string]string{}, ErrUserNotFound
	}
	redisMapRes := ur.RDB.HGetAll(context.Background(), fmt.Sprintf("user:%s", id))
	if redisMapRes.Err() != nil {
		ur.Logger.Error(logging.Redis, logging.Get, redisMapRes.Err().Error(), nil)
		return map[string]string{}, redisMapRes.Err()
	}
	ur.Logger.Info(logging.Redis, logging.Get, "", nil)
	return redisMapRes.Val(), nil
}
func (ur *UserRepo) Create(firstname, lastname, biography, username, password, email, phonenumber string) (models.User, *gorm.DB, error) {
	var u models.User
	u.Firstname = firstname
	u.Lastname = lastname
	u.Biography = biography
	u.Username = username
	u.Email = email
	u.PhoneNumber = phonenumber
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		ur.Logger.Error(logging.Internal, logging.HashPassword, err.Error(), nil)
		return u, nil, err
	}
	u.Password = hashedPassword
	tx := NewTx(ur.DB)
	err = tx.Create(&u).Error
	if err != nil {
		if utils.CheckErrorForWord(err, "email") {
			ur.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
			return u, nil, ErrEmailAlreadyExits
		} else if utils.CheckErrorForWord(err, "username") {
			ur.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
			return u, nil, ErrUsernameAlreadyExits
		} else if utils.CheckErrorForWord(err, "phone_number") {
			ur.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
			return u, nil, ErrPhoneNumberAlreadyExits
		} else {
			ur.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
			return u, nil, err
		}
	}
	err = ur.CreateChache(u)
	if err != nil {
		tx.Rollback()
		ur.Logger.Error(logging.Redis, logging.Set, err.Error(), nil)
		ur.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return u, nil, err
	}
	txj := tx.Commit()
	ur.Logger.Info(logging.Mysql, logging.Insert, "", nil)
	//retrun tx for rollback if jwt token can not be set
	return u, txj, nil
}
func (ur *UserRepo) UpdateById(id, firstname, lastname, biography, username string) (models.User, error) {
	var u models.User
	tx := NewTx(ur.DB)
	err := tx.First(&u, id).Error
	if err != nil {
		ur.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return u, err
	}
	u.Firstname = firstname
	u.Lastname = lastname
	u.Biography = biography
	u.Username = username
	err = tx.Save(&u).Error
	if err != nil {
		tx.Rollback()
		ur.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return u, err
	}
	err = ur.CreateChache(u)
	if err != nil {
		tx.Rollback()
		ur.Logger.Error(logging.Redis, logging.Set, err.Error(), nil)
		ur.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return u, err
	}
	ur.Logger.Info(logging.Mysql, logging.Update, "", nil)
	return u, nil
}
func (ur *UserRepo) DeleteById(id string) error {
	var u models.User
	tx := NewTx(ur.DB)
	err := tx.Delete(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
			return ErrUserNotFound
		}
	}
	id = strconv.Itoa(int(u.Id))
	err = ur.DeleteChacheById(id)
	if err != nil {
		tx.Rollback()
		ur.Logger.Error(logging.Redis, logging.Delete, err.Error(), nil)
		ur.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return err
	}
	ur.Logger.Info(logging.Mysql, logging.Delete, "", nil)
	return nil
}
func (ur *UserRepo) Verify(username, password string) (models.User, error) {
	var u models.User
	err := ur.DB.First(&u, "username=?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ur.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
			return u, ErrUserNotFound
		}
		ur.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return u, err
	}
	err = utils.CheckPassword(password, u.Password)
	if err != nil {
		if utils.CheckErrorForWord(err, "crypto/bcrypt") {
			ur.Logger.Error(logging.Mysql, logging.Verify, err.Error(), nil)
			return u, ErrUsernameOrPasswordWrong
		}
		ur.Logger.Error(logging.Mysql, logging.Verify, err.Error(), nil)
		return u, err
	}
	err = ur.CreateChache(u)
	if err != nil {
		ur.Logger.Info(logging.Mysql, logging.Verify, err.Error(), nil)
		return u, err
	}
	ur.Logger.Info(logging.Mysql, logging.Verify, "", nil)
	return u, nil
}

func (ur *UserRepo) CreateChache(u models.User) error {
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
	return redisRes.Err()
}
func (ur *UserRepo) DeleteChacheById(id string) error {
	redisRes := database.Rdb.Del(context.Background(), fmt.Sprintf("user:%s", id))
	return redisRes.Err()
}
func (ur *UserRepo) GetUsernameById(id string) (string, error) {
	user, err := ur.GetById(id)
	if err != nil {
		ur.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
		return "", err
	}
	return user["username"], err
}
