package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"loan-tracker-api/internal/domain/models"
)


func GenerateVerificationToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}


func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %v", err)
	}
	return string(hashedPassword), nil
}


func ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateResetToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}

	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}

func UpdateUserVerificationStatus(user *models.User) error {
	if user.IsActive {
		return errors.New("email already verified")
	}

	// Mark the user's email as verified
	user.IsActive = true

	return nil
}