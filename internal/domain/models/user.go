package models

import (
	"time"
	"github.com/google/uuid"
)

type User struct {
	ID                string    `bson:"_id,omitempty"`
	Email             string    `bson:"email"`
	Password          string    `bson:"password"`
	Name              string    `bson:"name"`
	Role              string    `bson:"role"`
	IsActive          bool      `bson:"is_active"`
	VerificationToken string    `bson:"verification_token"`
	CreatedAt         time.Time `bson:"created_at"`
	UpdatedAt         time.Time `bson:"updated_at"`
}


func NewUser(email, name, role string) (*User, error) {

	user := &User{
		ID:        generateID(),
		Email:     email,
		Name:      name,
		Role:      role,
		IsActive:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return user, nil
}


func (u *User) SetVerificationToken(token string) {
	u.VerificationToken = token
	u.UpdatedAt = time.Now() // Update the updated_at field
}

// Activate marks the user as active
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// Deactivate marks the user as inactive
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// generateID generates a unique ID for the user
func generateID() string {
	return uuid.NewString()
}
