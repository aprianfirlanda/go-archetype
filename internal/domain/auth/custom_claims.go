package auth

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	Roles []string `json:"roles"`
	jwt.RegisteredClaims
}
