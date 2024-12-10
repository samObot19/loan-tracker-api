package usecase

import (
	"context"
	"errors"
	"loan-tracker-api/internal/domain/models"
	"loan-tracker-api/internal/repository"
	e "loan-tracker-api/internal/infrastructure/email"
	"loan-tracker-api/pkg/utils"
	"loan-tracker-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserUsecase interface {
	RegisterUser(ctx context.Context, email, password, name, role string) (*models.User, error)
	VerifyEmail(ctx context.Context, token string) error
	LoginUser(ctx context.Context, email, password string) (*models.User, string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error) // New method for refreshing tokens
	GetUserProfile(ctx context.Context, userID string) (*models.User, error)
	RequestPasswordReset(ctx context.Context, email string) error
	UpdatePassword(ctx context.Context, token, newPassword string) error
}

type userUsecaseImpl struct {
	userRepo repository.UserRepository
}

func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecaseImpl{
		userRepo: userRepo,
	}
}

func (u *userUsecaseImpl) RegisterUser(ctx context.Context, email, password, name, role string) (*models.User, error) {
	if role != models.RoleUser && role != models.RoleAdmin {
		return nil, errors.New("invalid role")
	}

	user, err := models.NewUser(email, name, role)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	token, err := utils.GenerateVerificationToken()
	if err != nil {
		return nil, err
	}
	user.SetVerificationToken(token)
	err = u.userRepo.RegisterUser(ctx, user)
	if err != nil {
		return nil, err
	}

	err = e.SendVerificationEmail(email, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecaseImpl) VerifyEmail(ctx context.Context, token string) error {
	// Find user by verification token
	user, err := u.userRepo.FindUserByVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid or expired token")
	}

	// Update user's email verification status
	err = utils.UpdateUserVerificationStatus(user)
	if err != nil {
		return err
	}

	// Save the updated user information
	err = u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) LoginUser(ctx context.Context, email, password string) (*models.User, string, string, error) {
	user, err := u.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	if !user.IsActive {
		return nil, "", "", errors.New("account not active")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", "", errors.New("invalid email or password")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (u *userUsecaseImpl) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", errors.New("invalid refresh token")
	}
	accessToken, err := jwt.GenerateAccessToken(claims.UserID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (u *userUsecaseImpl) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	user, err := u.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *userUsecaseImpl) UpdatePassword(ctx context.Context, token, newPassword string) error {
	userID, err := jwt.VerifyResetToken(token)
	if err != nil {
		return err
	}

	user, err := u.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	err = u.userRepo.UpdateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *userUsecaseImpl) RequestPasswordReset(ctx context.Context, emailAddress string) error {
	// Step 1: Find user by email
	user, err := u.userRepo.FindUserByEmail(ctx, emailAddress)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("user not found")
	}

	// Step 2: Generate a password reset token
	resetToken, err := utils.GenerateResetToken()
	if err != nil {
		return err
	}

	// Step 3: Store the reset token with an expiration time
	expiration := time.Now().Add(1 * time.Hour) // Token expires in 1 hour
	err = u.userRepo.StorePasswordResetToken(ctx, user.ID, resetToken, expiration)
	if err != nil {
		return err
	}

	// Step 4: Send the password reset email
	resetLink := "https://yourdomain.com/reset-password?token=" + resetToken
	err = e.SendPasswordResetEmail(user.Email, resetLink)
	if err != nil {
		return err
	}

	// Step 5: Log the request
	utils.LogInfo("Password reset request for email: " + emailAddress)

	return nil
}
