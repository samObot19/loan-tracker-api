package validator

import (
	"errors"
	"regexp"
)


func ValidateEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	if !re.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

// validatePassword checks if the password meets criteria.
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	// Add more password strength validation if needed
	return nil
}

// validateRole checks if the role is valid.
func ValidateRole(role string) error {
	validRoles := []string{"user", "admin"}
	for _, r := range validRoles {
		if r == role {
			return nil
		}
	}
	return errors.New("invalid role")
}
