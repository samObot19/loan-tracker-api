package usecase

import (
	"context"
	"loan-tracker-api/internal/repository"
	"loan-tracker-api/internal/domain/models"
)

type AdminUsecase interface {
	ViewAllUsers(ctx context.Context) ([]*models.User, error) // Method to retrieve all users
	DeleteUser(ctx context.Context, userID string) error     // Method to delete a user by ID
}


type adminUsecaseImpl struct {
	userRepo repository.UserRepository
}

// NewAdminUsecase creates a new instance of AdminUsecase.
func NewAdminUsecase(userRepo repository.UserRepository) AdminUsecase {
	return &adminUsecaseImpl{
		userRepo: userRepo,
	}
}

// ViewAllUsers retrieves all users from the repository.
func (a *adminUsecaseImpl) ViewAllUsers(ctx context.Context) ([]*models.User, error) {
	return a.userRepo.FindAllUsers(ctx)
}

// DeleteUser deletes a user by ID from the repository.
func (a *adminUsecaseImpl) DeleteUser(ctx context.Context, userID string) error {
	return a.userRepo.DeleteUserByID(ctx, userID)
}
