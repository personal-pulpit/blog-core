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
type userManager struct {
	userPostgresRepo  repository.UserPostgresRepository
	authPostgresRepo repository.AuthPostgresRepository
}

func NewUserService(userPostgresRepo repository.UserPostgresRepository, authPostgresRepo repository.AuthPostgresRepository) UserService {
	return &userManager{
		userPostgresRepo:  userPostgresRepo,
		authPostgresRepo: authPostgresRepo,
	}
}
func (u *userManager) GetUserProfile(ID model.ID)(*model.User,error){
	userModel,err := u.userPostgresRepo.GetUserByID(ID)
	if err != nil{
		return nil,ErrNotFound
	}
	return userModel,nil
}
func (u *userManager) UpdateProfile(ID model.ID, FirstName, lastName, biography string) (*model.User,error) {
	user, err := u.userPostgresRepo.GetUserByID(ID)
	if err != nil {
		return nil,ErrNotFound
	}
	user.FirstName = FirstName
	user.LastName = lastName
	user.Biography = biography

	userModel, err := u.userPostgresRepo.UpdateByID(ID, user.FirstName, user.LastName, user.Biography)
	if err != nil {
		return nil,ErrUpdateUser
	}
	return userModel,nil
}

func (u *userManager) DeleteAccount(ID model.ID, password string) error {
	auth, err := u.authPostgresRepo.GetUserAuth(ID)
	if err != nil {
		return ErrNotFound
	}
	if !auth.EmailVerified {
		return ErrDeleteUser
	}
	err = u.authPostgresRepo.DeleteByID(ID)
	if err != nil{
		return ErrDeleteUser
	}
	err = u.userPostgresRepo.DeleteByID(ID)
	if err != nil {
		return ErrDeleteUser
	}
	return nil
}
func(u *userManager)Logout(token string)error{
	panic("not impl")
}