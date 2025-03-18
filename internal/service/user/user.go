package user

import (
	"blog/internal/model"
	"blog/internal/repository"
	"context"
)

type UserService interface {
	GetUserProfile(ctx context.Context, ID uint) (*model.User, error)
	UpdateProfile(ctx context.Context, ID uint, FirstName, lastName, biography string) (*model.User, error)
	DeleteAccount(ctx context.Context, ID uint, password string) error
}
type userManager struct {
	userPostgresRepo repository.UserRepository
	authPostgresRepo repository.AuthRepository
}

func NewUserService(userPostgresRepo repository.UserRepository, authPostgresRepo repository.AuthRepository) UserService {
	return &userManager{
		userPostgresRepo: userPostgresRepo,
		authPostgresRepo: authPostgresRepo,
	}
}
func (u *userManager) GetUserProfile(ctx context.Context, ID uint) (*model.User, error) {
	userModel, err := u.userPostgresRepo.GetUserByID(ctx, ID)
	if err != nil {
		return nil, ErrNotFound
	}
	return userModel, nil
}
func (u *userManager) UpdateProfile(ctx context.Context, ID uint, FirstName, lastName, biography string) (*model.User, error) {
	user, err := u.userPostgresRepo.GetUserByID(ctx, ID)
	if err != nil {
		return nil, ErrNotFound
	}
	user.FirstName = FirstName
	user.LastName = lastName
	user.Biography = biography

	userModel, err := u.userPostgresRepo.UpdateByID(ctx, ID, user.FirstName, user.LastName, user.Biography)
	if err != nil {
		return nil, ErrUpdateUser
	}
	return userModel, nil
}

func (u *userManager) DeleteAccount(ctx context.Context, ID uint, password string) error {
	auth, err := u.authPostgresRepo.GetUserAuth(ctx, ID)
	if err != nil {
		return ErrNotFound
	}

	if !auth.EmailVerified {
		return ErrDeleteUser
	}

	err = u.authPostgresRepo.DeleteByID(ctx, ID)
	if err != nil {
		return ErrDeleteUser
	}

	err = u.userPostgresRepo.DeleteByID(ctx, ID)
	if err != nil {
		return ErrDeleteUser
	}

	return nil
}
