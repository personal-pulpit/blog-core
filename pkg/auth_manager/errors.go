package auth_manager

import "errors"
import auth_manager "github.com/tahadostifam/go-auth-manager"

var (
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = auth_manager.ErrTokenExpired
)
