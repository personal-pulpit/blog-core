package repository

import "errors"

var (
	ErrEmailAlreadyExits    = errors.New("email already exits")
	ErrAuthNotFound         = errors.New("auth not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrLikeNotFound         = errors.New("like not found")
	ErrBookmarkNotFound     = errors.New("bookmark not found")
	ErrUserImageURLNotFound = errors.New("user image url not found")
	ErrCommentNotFound      = errors.New("comment not found")
	ErrArticleNotFound      = errors.New("article not found")
	ErrCategoryNotFound     = errors.New("category not found")
	ErrDatabase             = errors.New("database error")
)
