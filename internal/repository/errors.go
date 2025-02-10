package repository

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrNotModified      = errors.New("not modified")
	ErrNotDeleted       = errors.New("not deleted")
	ErrEmailAlreadyExits       = errors.New("email already exits")
)
