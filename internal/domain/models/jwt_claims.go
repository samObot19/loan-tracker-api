package models

import "github.com/golang-jwt/jwt/v4"


const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
