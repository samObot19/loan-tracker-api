package handler

import (
	"io/ioutil"
	"net/http"
	"loan-tracker-api/internal/usecase"
	"github.com/gin-gonic/gin"
	//"loan-tracker-api/pkg/utils"
)

type adminHandlerImpl struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminHandler(adminUsecase usecase.AdminUsecase) adminHandlerImpl {
	return adminHandlerImpl{
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

func (a *adminHandlerImpl) ViewSystemLogs(w http.ResponseWriter, r *http.Request) {
	logFile := "path/to/your/log/file.log" // Retrieve this path from config if necessary

	content, err := ioutil.ReadFile(logFile)
	if err != nil {
		http.Error(w, "Could not read log file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}
