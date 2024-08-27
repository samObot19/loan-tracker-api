package repository

import (
	"context"
	"loan-tracker-api/internal/domain/models"
	"loan-tracker-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
   
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type loanRepositoryImpl struct {
	collection *mongo.Collection
}

func NewLoanRepository(db *mongo.Database) repository.LoanRepository {
	return &loanRepositoryImpl{collection: db.Collection("loans")}
}

func (r *loanRepositoryImpl) CreateLoan(ctx context.Context, loan *models.Loan) error {
	_, err := r.collection.InsertOne(ctx, loan)
	return err
}

func (r *loanRepositoryImpl) FindLoanByID(ctx context.Context, id string) (*models.Loan, error) {
	var loan models.Loan
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&loan)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // No loan found
		}
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepositoryImpl) FindAllLoans(ctx context.Context, status models.LoanStatus, sortOrder int) ([]*models.Loan, error) {
	filter := bson.M{}
	if status != models.StatusAll {
		filter["status"] = status
	}

	opts := options.Find().SetSort(bson.D{{"created_at", sortOrder}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var loans []*models.Loan
	for cursor.Next(ctx) {
		var loan models.Loan
		if err := cursor.Decode(&loan); err != nil {
			return nil, err
		}
		loans = append(loans, &loan)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return loans, nil
}

func (r *loanRepositoryImpl) UpdateLoan(ctx context.Context, loan *models.Loan) error {
	filter := bson.M{"_id": loan.ID}
	update := bson.M{
		"$set": bson.M{
			"status":     loan.Status,
			"updated_at": loan.UpdatedAt,
		},
	}
	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *loanRepositoryImpl) DeleteLoan(ctx context.Context, loanID string) error {
	filter := bson.M{"_id": loanID}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
