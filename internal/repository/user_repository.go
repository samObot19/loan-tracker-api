package repository

import (
	"context"
	"loan-tracker-api/internal/domain/models"
	"time"
)


type UserRepository interface {
	RegisterUser(ctx context.Context, user *models.User) error
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByVerificationToken(ctx context.Context, token string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	FindAllUsers(ctx context.Context) ([]*models.User, error)
	DeleteUserByID(ctx context.Context, userID string) error
	StorePasswordResetToken(ctx context.Context, userID, token string, expiration time.Time) error
}



