package handler

import (
	"net/http"
	"loan-tracker-api/internal/usecase"
	"github.com/gin-gonic/gin"
)

type adminHandlerImpl struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase usecase.AdminUsecase) AdminHandler {
	return &adminHandlerImpl{
		adminUsecase: adminUsecase,
	}
}


func (h *adminHandlerImpl) ViewAllUsers(c *gin.Context) {
	users, err := h.adminUsecase.ViewAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}


func (h *adminHandlerImpl) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	err := h.adminUsecase.DeleteUser(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
