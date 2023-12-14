package repository

import (
	"context"
	"errors"
	"log"
	"golang.org/x/crypto/bcrypt"


	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sandlayth/supplier-api/public/model"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoRepository) create(user *model.User) (string, error) {
	result, err := r.collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return "", err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error converting inserted ID to primitive.ObjectID\n")
		return "", errors.New("inserted ID is not a primitive.ObjectID")
	}
	return insertedID.Hex(), nil
}

func (r *MongoRepository) findByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("User not found for email: %s\n", email)
			return nil, nil // Return nil, nil when user is not found
		}

		log.Printf("Error finding user by email: %v\n", err)
		return nil, err
	}

	log.Printf("User found for email: %s, ID: %s\n", email, user.ID)
	return user, nil
}

func (r *MongoRepository) findById(id string) (*model.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting ID to primitive.ObjectID: %v\n", err)
		return nil, err
	}

	user := &model.User{}
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(user)
	if err != nil {
		log.Printf("Error finding user by ID: %v\n", err)
		return nil, err
	}
	return user, nil
}

func (r *MongoRepository) update(user *model.User) error {
	_, err := r.collection.ReplaceOne(context.Background(), bson.M{"_id": user.ID}, user)
	if err != nil {
		log.Printf("Error updating user: %v\n", err)
		return err
	}
	return nil
}

func (r *MongoRepository) delete(id string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil
}

func (r *MongoRepository) listAll() (*[]model.User, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error listing all users: %v\n", err)
		return nil, err
	}

	var results []model.User

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Printf("Error on the cursor when listing all users: %v\n", err)
		return nil, err
	}

	return &results, nil
}

func (r *MongoRepository) Register(email, password, firstName, lastName string) (string, error) {
	// Check if the user with the given email already exists
	existingUser, err := r.findByEmail(email)
	if err == nil && existingUser != nil {
		log.Print(err)
		return "", errors.New("user with the given email already exists")
	}

	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	newUser := &model.User{
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
	}

	userID, err := r.create(newUser)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *MongoRepository) Login(email, password string) (string, error) {
	// Find user by email
	user, err := r.findByEmail(email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Compare hashed password with the provided password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate and return an authentication token (@TODO later use JWT token)
	// For simplicity, let's return the user ID as the token
	return user.ID.Hex(), nil
}

func (r *MongoRepository) Logout(userID string) error {
	// Perform any necessary cleanup or token invalidation logic
	// For simplicity, no additional logic is added
	return nil
}

func (r *MongoRepository) GetUserInfo(userID string) (*model.User, error) {
	// Find user by ID
	user, err := r.findById(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Omit sensitive information before returning user data
	user.Password = "" // Do not expose the hashed password

	return user, nil
}

func (r *MongoRepository) ListAll() (*[]model.User, error) {
	// Find user by ID
	users, err := r.listAll()
	if err != nil {
		return nil, err
	}

	return users, nil
}
