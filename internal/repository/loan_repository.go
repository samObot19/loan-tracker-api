package repository

import (
    "context"
    "loan-tracker-api/internal/domain/models"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *models.Loan) error
	FindLoanByID(ctx context.Context, id string) (*models.Loan, error)
	FindAllLoans(ctx context.Context, status models.LoanStatus, sortOrder int) ([]*models.Loan, error)
	UpdateLoan(ctx context.Context, loan *models.Loan) error
	DeleteLoan(ctx context.Context, loanID string) error
}