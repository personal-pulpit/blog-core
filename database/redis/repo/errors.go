package redis_repository

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrArticleNotFound = errors.New("article not found")
)