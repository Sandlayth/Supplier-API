package repository

import (
	"context"
	"errors"
	"net/mail"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/sandlayth/supplier-api/helper"
	"github.com/sandlayth/supplier-api/model"
)

type UserMongoRepository struct {
	collection *mongo.Collection
}

func NewUserMongoRepository(db *mongo.Database) *UserMongoRepository {
	return &UserMongoRepository{
		collection: db.Collection("users"),
	}
}

// CreateUser adds a new user to the database.
func (r *UserMongoRepository) CreateUser(user *model.User) error {
	if err := r.validateUser(user); err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	result, err := r.collection.InsertOne(context.Background(), user)
	// Update the purchase with the new ID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("inserted ID is not a primitive.ObjectID")
	}
	user.ID = insertedID
	return err
}

// GetUserByID retrieves a user by ID from the database.
func (r *UserMongoRepository) GetUserByID(id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email from the database.
func (r *UserMongoRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates an existing user in the database.
func (r *UserMongoRepository) UpdateUser(id string, updatedUser *model.User) error {
	if err := r.validateUser(updatedUser); err != nil {
		return err
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.MinCost)
	if err != nil {
		return err
	}
	updatedUser.Password = string(hashedPassword)
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updatedUser})
	updatedUser.ID = objectID
	return err
}

// DeleteUser removes a user from the database by ID.
func (r *UserMongoRepository) DeleteUser(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

// ListUsers retrieves a list of all users from the database.
func (r *UserMongoRepository) ListUsers() ([]model.User, error) {
	var users []model.User
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (r *UserMongoRepository) ListAll() (*[]model.User, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	var results []model.User

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return &results, nil
}

func (r *UserMongoRepository) GetTokens(user *model.User) (string, string, error) {
	// Generate refresh token
	refreshToken, err := helper.GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	// Generate access token
	accessToken, err := helper.GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}
	return refreshToken, accessToken, nil
}

func (r *UserMongoRepository) RenewTokens(userID string, refreshToken string) (string, string, error) {
	// Get the user from the DB
	dbUser, err := r.GetUserByID(userID)
	if err != nil {
		return "", "", err
	}
	// Verify the refresh token and extract claims
	claims, needsRefresh, err := helper.VerifyToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	if userID != claims.UserID.Hex() {
		return "", "", errors.New("invalid user ID in refresh token")
	}
	if needsRefresh {
		return r.GetTokens(dbUser)
	}
	// Return the original refresh token
	accessToken, err := helper.GenerateAccessToken(dbUser)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil

}

func (r *UserMongoRepository) ValidateUserCredentials(user *model.User) error {
	dbUser, err := r.GetUserByEmail(user.Email)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
}

func (r *UserMongoRepository) validateUser(user *model.User) error {
	err := "Error when validating user input: "
	if _, e := mail.ParseAddress(user.Email); e != nil {
		return errors.New(err + "invalid email")
	}
	if len(strings.TrimSpace(user.FirstName)) == 0 {
		return errors.New(err + "invalid firstName field")
	}
	if len(strings.TrimSpace(user.LastName)) == 0 {
		return errors.New(err + "invalid lastName field")
	}
	if len(strings.TrimSpace(user.Password)) == 0 {
		return errors.New(err + "invalid password field")
	}
	if user.Role != "manager" && user.Role != "admin" {
		return errors.New(err + "invalid role field")
	}
	return nil
}
