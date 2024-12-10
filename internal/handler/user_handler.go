package handler

import (
	"net/http"
	"strings"
	"context"
	"loan-tracker-api/internal/usecase"
	"loan-tracker-api/pkg/utils"
	"loan-tracker-api/pkg/jwt"
	"loan-tracker-api/internal/infrastructure/email"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap" // Assuming you are using zap for logging
)

type userHandlerImpl struct {
	userUsecase usecase.UserUsecase
	logger      *zap.Logger
}

func NewUserHandler(userUsecase usecase.UserUsecase, logger *zap.Logger)  userHandlerImpl {
	return userHandlerImpl{
		userUsecase: userUsecase,
		logger:      logger,
	}
}

func (h *userHandlerImpl) RegisterUser(c *gin.Context) {
	var reqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Role     string `json:"role"`
	}

	if err := c.BindJSON(&reqBody); err != nil {
		h.logger.Error("Failed to bind JSON for registration", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateVerificationToken()
	if err != nil {
		h.logger.Error("Failed to generate verification token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate verification token"})
		return
	}

	err = email.SendVerificationEmail(reqBody.Email, token)
	if err != nil {
		h.logger.Error("Failed to send verification email", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send verification email"})
		return
	}

	_, err= h.userUsecase.RegisterUser(context.Background(), reqBody.Email, reqBody.Password, reqBody.Name, reqBody.Role)
	if err != nil {
		h.logger.Error("Failed to register user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("User registered successfully", zap.String("email", reqBody.Email))
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully. Please check your email to verify your account."})
}

func (h *userHandlerImpl) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	email := c.Query("email")

	if token == "" || email == "" {
		h.logger.Warn("Email verification failed due to missing parameters", zap.String("token", token), zap.String("email", email))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token and email are required"})
		return
	}

	err := h.userUsecase.VerifyEmail(c, token)
	if err != nil {
		h.logger.Error("Failed to verify email", zap.String("email", email), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Email verified successfully", zap.String("email", email))
	c.JSON(http.StatusOK, gin.H{"message": "Email successfully verified"})
}

func (h *userHandlerImpl) RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind JSON for refresh token", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.userUsecase.RefreshToken(c, request.RefreshToken)
	if err != nil {
		h.logger.Error("Failed to refresh token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Token refreshed successfully")
	c.JSON(http.StatusOK, gin.H{"access_token": accessToken})
}

func (h *userHandlerImpl) GetUserProfile(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		h.logger.Warn("Authorization header missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims, err := jwt.ParseToken(tokenString)
	if err != nil {
		h.logger.Error("Invalid or expired token", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	userID := claims.UserID

	user, err := h.userUsecase.GetUserProfile(c, userID)
	if err != nil {
		h.logger.Error("Failed to get user profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Retrieved user profile successfully", zap.String("userID", userID))
	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"email": user.Email,
		"name":  user.Name,
		"role":  user.Role,
	})
}

func (h *userHandlerImpl) UpdatePassword(c *gin.Context) {
	var request struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=8"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind JSON for password update", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.UpdatePassword(c, request.Token, request.NewPassword)
	if err != nil {
		h.logger.Error("Failed to update password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Password updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func (h *userHandlerImpl) RequestPasswordReset(c *gin.Context) {
	var request struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Error("Failed to bind JSON for password reset request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userUsecase.RequestPasswordReset(c, request.Email)
	if err != nil {
		h.logger.Error("Failed to request password reset", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Password reset link sent", zap.String("email", request.Email))
	c.JSON(http.StatusOK, gin.H{"message": "Password reset link sent"})
}
