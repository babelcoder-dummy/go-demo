package auth

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Role string `json:"role"`
	Type string `json:"typ"`
	jwt.RegisteredClaims
}
