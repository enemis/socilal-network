package auth

import (
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	TokenId string `json:"token_id"`
	UserID  string `json:"user_id"`
}

func (c Claims) Valid() error {
	return c.StandardClaims.Valid()
}
