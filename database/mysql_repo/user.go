package mysql_repository

import (
	"blog/database"
	"blog/internal/model"
	"blog/internal/repository"
	"blog/pkg/logging"
	"blog/utils"
	"errors"

	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type userRepo struct {
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

func NewUserRepository() repository.UserRepository {
	return &userRepo{
		DB:     database.GetMysqlDB(),
		RDB:    database.GetRedisDB(),
		Logger: logging.MyLogger,
	}
}
func (u *userRepo) GetAll() ([]map[string]string, error) {
	var users []map[string]string
	keys, err := u.RDB.Keys(context.Background(), "user:*").Result()
	if err != nil {
		u.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
		return users, err
	}
	for _, key := range keys {
		userMap, err := u.RDB.HGetAll(context.Background(), key).Result()
		if err != nil {
			u.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
			return []map[string]string{}, err
		}
		users = append(users, userMap)
	}
	u.Logger.Info(logging.Redis, logging.Get, "", nil)
	return users, nil
}
func (u *userRepo) GetByID(ID string) (map[string]string, error) {
	exists := u.RDB.Exists(context.Background(), fmt.Sprintf("user:%s", ID))
	if exists.Val() == 0 {
		u.Logger.Error(logging.Redis, logging.Get, ErrUserNotFound.Error(), nil)
		return map[string]string{}, ErrUserNotFound
	}
	redisMapRes := u.RDB.HGetAll(context.Background(), fmt.Sprintf("user:%s", ID))
	if redisMapRes.Err() != nil {
		u.Logger.Error(logging.Redis, logging.Get, redisMapRes.Err().Error(), nil)
		return map[string]string{}, redisMapRes.Err()
	}
	u.Logger.Info(logging.Redis, logging.Get, "", nil)
	return redisMapRes.Val(), nil
}
func (u *userRepo) Create(firstname, lastname, biography, username, password, email, phonenumber string) (model.User, *gorm.DB, error) {
	var user model.User
	user.Firstname = firstname
	user.Lastname = lastname
	user.Biography = biography
	user.Username = username
	user.Email = email
	user.PhoneNumber = phonenumber
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		u.Logger.Error(logging.Internal, logging.HashPassword, err.Error(), nil)
		return user, nil, err
	}
	user.Password = hashedPassword
	tx := NewTx(u.DB)
	err = tx.Create(&u).Error
	if err != nil {
		if utils.CheckErrorForWord(err, "email") {
			u.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
			return user, nil, ErrEmailAlreadyExits
		} else if utils.CheckErrorForWord(err, "username") {
			u.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
			return user, nil, ErrUsernameAlreadyExits
		} else if utils.CheckErrorForWord(err, "phone_number") {
			u.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
			return user, nil, ErrPhoneNumberAlreadyExits
		} else {
			u.Logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
			return user, nil, err
		}
	}
	err = u.CreateChache(user)
	if err != nil {
		tx.Rollback()
		u.Logger.Error(logging.Redis, logging.Set, err.Error(), nil)
		u.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return user, nil, err
	}
	txj := tx.Commit()
	u.Logger.Info(logging.Mysql, logging.Insert, "", nil)
	//retrun tx for rollback if jwt token can not be set
	return user, txj, nil
}
func (u *userRepo) UpdateByID(ID, firstname, lastname, biography, username string) (model.User, error) {
	var user model.User
	tx := NewTx(u.DB)
	err := tx.First(&u, ID).Error
	if err != nil {
		u.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return user, err
	}
	user.Firstname = firstname
	user.Lastname = lastname
	user.Biography = biography
	user.Username = username
	err = tx.Save(&u).Error
	if err != nil {
		tx.Rollback()
		u.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return user, err
	}
	err = u.CreateChache(user)
	if err != nil {
		tx.Rollback()
		u.Logger.Error(logging.Redis, logging.Set, err.Error(), nil)
		u.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return user, err
	}
	u.Logger.Info(logging.Mysql, logging.Update, "", nil)
	return user, nil
}
func (u *userRepo) DeleteByID(ID string) error {
	var user model.User
	tx := NewTx(u.DB)
	err := tx.Delete(&u, ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
			return ErrUserNotFound
		}
	}
	ID = strconv.Itoa(int(user.ID))
	err = u.DeleteChacheById(ID)
	if err != nil {
		tx.Rollback()
		u.Logger.Error(logging.Redis, logging.Delete, err.Error(), nil)
		u.Logger.Error(logging.Mysql, logging.Rollback, err.Error(), nil)
		return err
	}
	u.Logger.Info(logging.Mysql, logging.Delete, "", nil)
	return nil
}
func (u *userRepo) Verify(username, password string) (model.User, error) {
	var user model.User
	err := u.DB.First(&u, "username=?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			u.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
			return user, ErrUserNotFound
		}
		u.Logger.Error(logging.Mysql, logging.Select, err.Error(), nil)
		return user, err
	}
	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		if utils.CheckErrorForWord(err, "crypto/bcrypt") {
			u.Logger.Error(logging.Mysql, logging.Verify, err.Error(), nil)
			return user, ErrUsernameOrPasswordWrong
		}
		u.Logger.Error(logging.Mysql, logging.Verify, err.Error(), nil)
		return user, err
	}
	err = u.CreateChache(user)
	if err != nil {
		u.Logger.Info(logging.Mysql, logging.Verify, err.Error(), nil)
		return user, err
	}
	u.Logger.Info(logging.Mysql, logging.Verify, "", nil)
	return user, nil
}

func (u *userRepo) CreateChache(user model.User) error {
	redisRes := u.RDB.HMSet(context.Background(), fmt.Sprintf("user:%d", user.ID), map[string]interface{}{
		"firstname":   user.Firstname,
		"lastname":    user.Lastname,
		"biography":   user.Biography,
		"username":    user.Username,
		"email":       user.Email,
		"phonenumber": user.PhoneNumber,
		"role":        user.Role,
		"createdAt":   user.CreatedAt,
		"updatedAt":   user.UpdatedAt,
	})
	return redisRes.Err()
}
func (u *userRepo) DeleteChacheById(ID string) error {
	redisRes := u.RDB.Del(context.Background(), fmt.Sprintf("user:%s", ID))
	return redisRes.Err()
}
func (u *userRepo) GetUsernameById(ID string) (string, error) {
	user, err := u.GetByID(ID)
	if err != nil {
		u.Logger.Error(logging.Redis, logging.Get, err.Error(), nil)
		return "", err
	}
	return user["username"], err
}
