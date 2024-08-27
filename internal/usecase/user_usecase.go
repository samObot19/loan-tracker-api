package usecase

import (
	"context"
	"errors"
	"loan-tracker-api/internal/domain/models"
	"loan-tracker-api/internal/repository"
	"loan-tracker-api/pkg/password"
	"loan-tracker-api/pkg/utils"
	"loan-tracker-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
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

	hashedPassword, err := password.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	token, err := utils.GenerateVerificationToken(email)
	if err != nil {
		return nil, err
	}
	user.SetVerificationToken(token)
	err = u.userRepo.RegisterUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// Send verification email
	err = utils.SendVerificationEmail(email, token)
	if err != nil {
		return nil, err
	}

	return user, nil
}



func (u *userUsecaseImpl) VerifyEmail(ctx context.Context, token string) error {
	user, err := u.userRepo.FindUserByVerificationToken(ctx, token)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("invalid or expired token")
	}

	user.VerifyEmail()
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

	accessToken, err := jwt.GenerateAccessToken(user.ID.Hex())
	if err != nil {
		return nil, "", "", err
	}

	refreshToken, err := jwt.GenerateRefreshToken(user.ID.Hex())
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

func (h *userHandlerImpl) RequestPasswordReset(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.RequestPasswordReset(c, request.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent"})
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