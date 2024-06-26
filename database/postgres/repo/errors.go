package postgres_repository

import "errors"

var (
	ErrEmailAlreadyExits       = errors.New("email already exits")
	ErrUsernameAlreadyExits    = errors.New("username already exits")
	ErrPhoneNumberAlreadyExits = errors.New("phone number already exits")
	ErrUserNotFound            = errors.New("user not found")
	ErrUsernameOrPasswordWrong = errors.New("username or password wrong")
)
