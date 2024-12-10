package handler

import (
    "net/http"
    "loan-tracker-api/internal/usecase"
    "loan-tracker-api/internal/domain/models"
    "github.com/gin-gonic/gin"
)

type LoanHandler struct {
    loanUsecase usecase.LoanUsecase
}

func NewLoanHandler(loanUsecase usecase.LoanUsecase) *LoanHandler {
    return &LoanHandler{loanUsecase: loanUsecase}
}

func (h *LoanHandler) ApplyForLoan(c *gin.Context) {
    var input struct {
        Amount float64 `json:"amount" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    userID, _ := c.Get("user_id") // Assuming user ID is set in the context by middleware
    loan, err := h.loanUsecase.ApplyForLoan(c, userID.(string), input.Amount)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to apply for loan"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"loan": loan})
}

func (h *LoanHandler) GetLoanStatus(c *gin.Context) {
    loanID := c.Param("id")
    userID, _ := c.Get("user_id") // Assuming user ID is set in the context by middleware

    loan, err := h.loanUsecase.GetLoanStatus(c, userID.(string), loanID)
    if err != nil {
        if err.Error() == "unauthorized access to this loan" {
            c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
        } else {
            c.JSON(http.StatusNotFound, gin.H{"error": "Loan not found"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"loan": loan})
}

func (h *LoanHandler) GetAllLoans(c *gin.Context) {
    status := c.Query("status")
    order := c.DefaultQuery("order", "asc")

    loans, err := h.loanUsecase.GetAllLoans(c.Request.Context(), status, order)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, loans)
}

func (h *LoanHandler) UpdateLoanStatus(c *gin.Context) {
    loanID := c.Param("id")
    statusStr := c.Query("status")

    var status models.LoanStatus
    switch statusStr {
    case "approved":
        status = models.StatusApproved
    case "rejected":
        status = models.StatusRejected
    default:
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
        return
    }

    loan, err := h.loanUsecase.UpdateLoanStatus(c.Request.Context(), loanID, status)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, loan)
}

func (h *LoanHandler) DeleteLoan(c *gin.Context) {
    loanID := c.Param("id")

    err := h.loanUsecase.DeleteLoan(c.Request.Context(), loanID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Loan deleted successfully"})
}