package auth_helper

import "errors"


var (
	ErrTokenIsInvalid = errors.New("token is invalid")
	ErrTokenUndefined = errors.New("token is undefind")
)