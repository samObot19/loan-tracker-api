package models

import "time"

type LoanStatus string

const (
    StatusPending  LoanStatus = "pending"
    StatusApproved LoanStatus = "approved"
    StatusRejected LoanStatus = "rejected"
	StatusAll      = "all"
)

type Loan struct {
    ID        string    `json:"id" bson:"_id"`
    UserID    string    `json:"user_id" bson:"user_id"`
    Amount    float64   `json:"amount" bson:"amount"`
    Status    LoanStatus `json:"status" bson:"status"`
    CreatedAt time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
