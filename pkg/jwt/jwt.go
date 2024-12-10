package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"loan-tracker-api/internal/domain/models"
)

var jwtSecret = []byte("your_secret_key")

// CustomClaims includes the standard JWT claims and custom claims (userID).
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// ParseToken parses and validates a JWT token string and returns the custom claims.
func ParseToken(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// GenerateToken generates a JWT token with custom claims.
func GenerateToken(claims *models.JWTClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// GenerateAccessToken generates an access token for a user.
func GenerateAccessToken(userID string) (string, error) {
	claims := &models.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // Access token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return GenerateToken(claims)
}

// GenerateRefreshToken generates a refresh token for a user.
func GenerateRefreshToken(userID string) (string, error) {
	claims := &models.JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // Refresh token expires in 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return GenerateToken(claims)
}

// ValidateToken validates a JWT token and returns the claims if valid.
func ValidateToken(tokenStr string) (*models.JWTClaims, error) {
	return ParseToken(tokenStr)
}

func ValidateRefreshToken(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

func VerifyResetToken(tokenString string) (string, error) {
	claims := &models.JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid or expired reset token")
	}

	return claims.UserID, nil
}