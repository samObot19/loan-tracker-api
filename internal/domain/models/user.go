package models

import (
	"time"
	"loan-tracker-api/pkg/validator"
)

// User represents a user in the system.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUser creates a new user instance after validating the input.
func NewUser(email, password, name, role string) (*User, error) {
	if err := validator.ValidateEmail(email); err != nil {
		return nil, err
	}

	if err := validator.ValidatePassword(password); err != nil {
		return nil, err
	}

	if err := validator.ValidateRole(role); err != nil {
		return nil, err
	}

	user := &User{
		ID:        generateID(),
		Email:     email,
		Password:  hashPassword(password),
		Name:      name,
		Role:      role,
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}

// Activate sets the user's status to active.
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate sets the user's status to inactive.
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Helper functions (assuming they are defined elsewhere)
func generateID() string {
	return "some-unique-id"
}

func hashPassword(password string) string {
	// Implementation of password hashing
	return "hashed-password"
}
