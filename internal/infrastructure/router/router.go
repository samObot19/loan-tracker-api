package router

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"loan-tracker-api/internal/handler"
	"loan-tracker-api/internal/repository/mongo"
	"loan-tracker-api/internal/repository"

	"loan-tracker-api/internal/infrastructure/middleware"
	"go.uber.org/zap"
)

func InitRoutes(r *gin.Engine, db *mongo.Client, logger *zap.Logger) {
	userRepo := repository.NewMongoUserRepository(db.Database("loan-tracker"))
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase, logger)

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("/register", userHandler.RegisterUser)
		userRoutes.GET("/verify-email", userHandler.VerifyEmail)
		userRoutes.POST("/login", userHandler.Login)
		userRoutes.POST("/token/refresh", userHandler.RefreshToken)

		userRoutes.Use(middleware.JWTAuthMiddleware())
		{
			userRoutes.GET("/profile", userHandler.GetUserProfile)
			userRoutes.POST("/password-reset", userHandler.RequestPasswordReset)
			// userRoutes.POST("/password-reset/confirm", userHandler.PasswordResetConfirm)
		}
	}

	adminRepo := repository.NewMongoAdminRepository(db.Database("loan-tracker")) // Assuming you have this
	adminUsecase := usecase.NewAdminUsecase(adminRepo)
	adminHandler := handler.NewAdminHandler(adminUsecase, logger)

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.JWTAuthorization(), middleware.AdminOnly())
	{
		adminRoutes.GET("/users", adminHandler.GetAllUsers)
		adminRoutes.DELETE("/users/:id", adminHandler.DeleteUser)
		adminRoutes.GET("/loans", adminHandler.ViewAllLoans)
		adminRoutes.PATCH("/loans/:id/status", adminHandler.UpdateLoanStatus)
		adminRoutes.DELETE("/loans/:id", adminHandler.DeleteLoan)
		adminRoutes.GET("/logs", adminHandler.ViewSystemLogs)
	}
}
