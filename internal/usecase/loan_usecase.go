package usecase

import (
    "context"
    "loan-tracker-api/internal/domain/models"
    "loan-tracker-api/internal/repository"
    "github.com/google/uuid"
	"loan-tracker-api/pkg/utils"
    "time"
	"fmt"
	"errors"
)

type LoanUsecase interface {
    ApplyForLoan(ctx context.Context, userID string, amount float64) (*models.Loan, error)
	GetLoanStatus(ctx context.Context, userID, loanID string) (*models.Loan, error)
	GetAllLoans(ctx context.Context, status string, order string) ([]*models.Loan, error)
	UpdateLoanStatus(ctx context.Context, loanID string, status models.LoanStatus) (*models.Loan, error)
	DeleteLoan(ctx context.Context, loanID string) error
}

type loanUsecaseImpl struct {
    loanRepo repository.LoanRepository
}

func NewLoanUsecase(loanRepo repository.LoanRepository) LoanUsecase {
    return &loanUsecaseImpl{loanRepo: loanRepo}
}

func (l *loanUsecaseImpl) ApplyForLoan(ctx context.Context, userID string, amount float64) (*models.Loan, error) {
    loan := &models.Loan{
        ID:        uuid.New().String(),
        UserID:    userID,
        Amount:    amount,
        Status:    models.StatusPending,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    err := l.loanRepo.CreateLoan(ctx, loan)
    if err != nil {
        return nil, err
    }
	utils.LogInfo("Loan application submitted by user " + userID + " for amount " + fmt.Sprintf("%.2f", amount))
    return loan, nil
}

func (l *loanUsecaseImpl) GetLoanStatus(ctx context.Context, userID, loanID string) (*models.Loan, error) {
    loan, err := l.loanRepo.FindLoanByID(ctx, loanID)
    if err != nil {
        return nil, err
    }

    if loan.UserID != userID {
        return nil, errors.New("unauthorized access to this loan")
    }

    return loan, nil
}

func (l *loanUsecaseImpl) GetAllLoans(ctx context.Context, status string, order string) ([]*models.Loan, error) {
    var loanStatus models.LoanStatus
    switch status {
    case "pending":
        loanStatus = models.StatusPending
    case "approved":
        loanStatus = models.StatusApproved
    case "rejected":
        loanStatus = models.StatusRejected
    default:
        loanStatus = models.StatusAll
    }

    sortOrder := 1 // Ascending
    if order == "desc" {
        sortOrder = -1
    }

    loans, err := l.loanRepo.FindAllLoans(ctx, loanStatus, sortOrder)
    if err != nil {
        return nil, err
    }

    return loans, nil
}

func (l *loanUsecaseImpl) UpdateLoanStatus(ctx context.Context, loanID string, status models.LoanStatus) (*models.Loan, error) {
	// Log the incoming request
	utils.LogInfo(fmt.Sprintf("Received request to update loan status for loan ID %s to %s", loanID, status))

	// Find the loan by ID
	loan, err := l.loanRepo.FindLoanByID(ctx, loanID)
	if err != nil {
		utils.LogError(fmt.Sprintf("Failed to find loan with ID %s: %v", loanID, err))
		return nil, err
	}

	loan.Status = status
	loan.UpdatedAt = time.Now()

	utils.LogInfo(fmt.Sprintf("Updating status of loan ID %s to %s", loanID, status))

	err = l.loanRepo.UpdateLoan(ctx, loan)
	if err != nil {
		utils.LogError(fmt.Sprintf("Failed to update loan status for loan ID %s: %v", loanID, err))
		return nil, err
	}

	utils.LogInfo(fmt.Sprintf("Successfully updated loan status for loan ID %s to %s", loanID, status))
	return loan, nil
}
func (l *loanUsecaseImpl) DeleteLoan(ctx context.Context, loanID string) error {
    //loan, err := l.loanRepo.FindLoanByID(ctx, loanID)
    // if err != nil {
    //     return err
    // }

    // Add any additional checks if needed (e.g., permissions)

    err := l.loanRepo.DeleteLoan(ctx, loanID)
    if err != nil {
        return err
    }

    return nil
}