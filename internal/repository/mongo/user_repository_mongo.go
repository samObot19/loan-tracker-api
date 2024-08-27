package repository

import (
	"context"
	"errors"
	"loan-tracker-api/internal/domain/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type MongoUserRepository struct {
	collection *mongo.Collection
}


func NewMongoUserRepository(db *mongo.Database) UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}


func (r *MongoUserRepository) RegisterUser(ctx context.Context, user *models.User) error {
	existingUser, _ := r.FindUserByEmail(ctx, user.Email)
	if existingUser != nil {
		return errors.New("user with this email already exists")
	}

	_, err := r.collection.InsertOne(ctx, user)
	return err
}


func (r *MongoUserRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}


func (r *MongoUserRepository) FindUserByVerificationToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"verification_token": token}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}


func (r *MongoUserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := r.collection.UpdateByID(ctx, user.ID, bson.M{"$set": user})
	return err
}
