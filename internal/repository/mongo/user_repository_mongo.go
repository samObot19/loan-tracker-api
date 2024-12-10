package repository

import (
	"context"
	"errors"
	"loan-tracker-api/internal/domain/models"
	"loan-tracker-api/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


type MongoUserRepository struct {
	collection *mongo.Collection
}


func NewMongoUserRepository(db *mongo.Database) repository.UserRepository {
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

func (r *MongoUserRepository) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	// Convert the userID string to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}

	// Create a filter to search by the user ID
	filter := bson.M{"_id": objectID}
	var user models.User

	// Find one document in the collection and decode it into the user model
	err = r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// No user was found with the given ID
			return nil, nil
		}
		// Another error occurred while finding the user
		return nil, err
	}

	// Return the found user
	return &user, nil
}


func (r *MongoUserRepository) FindAllUsers(ctx context.Context) ([]*models.User, error) {
	var users []*models.User
	filter := bson.M{}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)


	for cursor.Next(ctx) {
		var user models.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *MongoUserRepository) DeleteUserByID(ctx context.Context, userID string) error {
	// Convert the userID string to a MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	// Create a filter to search by the user ID
	filter := bson.M{"_id": objectID}

	// Delete the document from the collection
	_, err = r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r *MongoUserRepository) StorePasswordResetToken(ctx context.Context, userID, token string, expiration time.Time) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$set": bson.M{
			"password_reset_token":  token,
			"password_reset_expiry": expiration,
		}},
	)
	return err
}