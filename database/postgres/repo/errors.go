package postgres_repository

import "errors"

var (
	ErrEmailAlreadyExits       = errors.New("email already exits")
	ErrUserNotFound            = errors.New("user not found")
	ErrEmailOrPasswordWrong = errors.New("username or password wrong")
)
