package repository

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sandlayth/supplier-api/public/entity"
)

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoRepository) Create(user *entity.User) (string, error) {
	userDAO := entity.NewUserDAO(user)
	result, err := r.collection.InsertOne(context.Background(), userDAO)
	if err != nil {
		log.Printf("Error creating user: %v\n", err)
		return "", err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error converting inserted ID to primitive.ObjectID\n")
		return "", errors.New("inserted ID is not a primitive.ObjectID")
	}
	return insertedID.String(), nil
}

func (r *MongoRepository) FindByEmail(email string) (*entity.User, error) {
	userDAO := &entity.UserDAO{}
	err := r.collection.FindOne(context.Background(), bson.M{"email": email}).Decode(userDAO)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("User not found for email: %s\n", email)
			return nil, nil // Return nil, nil when user is not found
		}

		log.Printf("Error finding user by email: %v\n", err)
		return nil, err
	}

	log.Printf("User found for email: %s, ID: %s\n", email, userDAO.ID.Hex())
	return userDAO.ToUser(), nil
}

func (r *MongoRepository) FindById(id string) (*entity.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Printf("Error converting ID to primitive.ObjectID: %v\n", err)
		return nil, err
	}

	userDAO := &entity.UserDAO{}
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(userDAO)
	if err != nil {
		log.Printf("Error finding user by ID: %v\n", err)
		return nil, err
	}
	return userDAO.ToUser(), nil
}

func (r *MongoRepository) Update(user *entity.User) error {
	userDAO := entity.NewUserDAO(user)
	_, err := r.collection.ReplaceOne(context.Background(), bson.M{"_id": userDAO.ID}, userDAO)
	if err != nil {
		log.Printf("Error updating user: %v\n", err)
		return err
	}
	return nil
}

func (r *MongoRepository) Delete(id string) error {
	_, err := r.collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		log.Printf("Error deleting user: %v\n", err)
		return err
	}
	return nil
}

func (r *MongoRepository) ListAll() (*[]entity.User, error) {
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Printf("Error listing all users: %v\n", err)
		return nil, err
	}

	var results []entity.User

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Printf("Error on the cursor when listing all users: %v\n", err)
		return nil, err
	}

	return &results, nil
}
