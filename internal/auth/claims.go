package auth

import "github.com/golang-jwt/jwt/v5"

//Custom claims

type Claims struct {
	UserID uint64 `json:"user_id"`
	jwt.RegisteredClaims
}
