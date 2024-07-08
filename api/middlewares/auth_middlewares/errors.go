package auth_middlewares

import "errors"

var (
	ErrYouAreUnAuthorized = errors.New("you are unauthorized")
	ErrSomeTimesWentWrong = errors.New("sometimes went wrong")
)
