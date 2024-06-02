package helpers

import (
	"blog/config"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

var (
	secret                     =[]byte(config.Cfg.Jwt.Secret)
	errInvalidToken            = errors.New("token is invalid")
	errUnexpectedSigningMethod = errors.New("unexpected signin method")
)

func CreateToken(Id uint) (string, error) {
        fmt.Println("secret:",secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       Id,
	})
	tokenString, err := token.SignedString(secret)
	return tokenString, err
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	tkn, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		//check
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errUnexpectedSigningMethod
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := tkn.Claims.(jwt.MapClaims); ok && tkn.Valid {
		return claims, nil
	}
	return nil, errInvalidToken
}
