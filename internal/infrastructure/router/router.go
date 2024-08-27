package router

import (
	"github.com/gin-gonic/gin"
	"loan-tracker-api/internal/handler"
	"loan-tracker-api/internal/usecase"
	"loan-tracker-api/internal/repository"
	"loan-tracker-api/internal/infrastructure/database"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Connect to MongoDB
	db, err := database.ConnectDB("mongodb://localhost:27017", "loan_tracker")
	if err != nil {
		panic(err)
	}

	// Initialize repositories
	userRepo := repository.NewMongoUserRepository(db)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo)
	adminUsecase := usecase.NewAdminUsecase(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUsecase)
	adminHandler := handler.NewAdminHandler(adminUsecase)

	// Define routes
	r.POST("/users/register", userHandler.RegisterUser)
	r.GET("/users/verify-email", userHandler.VerifyEmail)
	r.POST("/users/login", userHandler.Login)
	r.POST("/users/token/refresh", userHandler.RefreshToken)
	r.GET("/users/profile", userHandler.GetUserProfile)
	r.POST("/users/password-reset", userHandler.RequestPasswordReset)
	r.POST("/users/password-update", userHandler.UpdatePassword)

	// Admin routes
	r.GET("/admin/users", adminHandler.ViewAllUsers)
	r.DELETE("/admin/users/:id", adminHandler.DeleteUser)

	return r
}
