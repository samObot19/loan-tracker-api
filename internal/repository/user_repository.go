package repository

import (
	"context"
	"loan-tracker-api/internal/domain/models"
)


type UserRepository interface {
	RegisterUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByVerificationToken(ctx context.Context, token string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

