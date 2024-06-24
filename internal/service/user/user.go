package user

import (
	"blog/internal/model"
	"blog/internal/repository"
)

type UserService interface {
	GetUserProfile(ID model.ID)(*model.User,error)
	UpdateProfile(ID model.ID, FirstName, lastName, biography string) (*model.User,error)
	DeleteAccount(ID model.ID, password string) error
	Logout(token string)error
}
type UserManager struct {
	userMysqRepo  repository.UserMysqlRepository
	authMysqlRepo repository.AuthMysqlRepository
}

func NewUserService(userMysqlRepo repository.UserMysqlRepository, authMysqlRepo repository.AuthMysqlRepository) UserService {
	return &UserManager{
		userMysqRepo:  userMysqlRepo,
		authMysqlRepo: authMysqlRepo,
	}
}
func (u *UserManager) GetUserProfile(ID model.ID)(*model.User,error){
	userModel,err := u.userMysqRepo.GetUserByID(ID)
	if err != nil{
		return nil,ErrNotFound
	}
	return userModel,nil
}
func (u *UserManager) UpdateProfile(ID model.ID, FirstName, lastName, biography string) (*model.User,error) {
	user, err := u.userMysqRepo.GetUserByID(ID)
	if err != nil {
		return nil,ErrNotFound
	}
	user.FirstName = FirstName
	user.LastName = lastName
	user.Biography = biography

	userModel, err := u.userMysqRepo.UpdateByID(ID, user.FirstName, user.LastName, user.Biography)
	if err != nil {
		return nil,ErrUpdateUser
	}
	return userModel,nil
}

func (u *UserManager) DeleteAccount(ID model.ID, password string) error {
	auth, err := u.authMysqlRepo.GetUserAuth(ID)
	if err != nil {
		return ErrNotFound
	}
	if !auth.EmailVerified {
		return ErrDeleteUser
	}
	err = u.authMysqlRepo.DeleteByID(ID)
	if err != nil{
		return ErrDeleteUser
	}
	err = u.userMysqRepo.DeleteByID(ID)
	if err != nil {
		return ErrDeleteUser
	}
	return nil
}
func(u *UserManager)Logout(token string)error{
	panic("not impl")
}