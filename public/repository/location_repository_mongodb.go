package repository

import (
	"context"

	"github.com/sandlayth/supplier-api/public/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationMongoRepository struct {
	collection *mongo.Collection
}

func NewLocationMongoRepository(db *mongo.Database) *LocationMongoRepository {
	return &LocationMongoRepository{
		collection: db.Collection("locations"),
	}
}

// CreateLocation adds a new location to the database.
func (r *LocationMongoRepository) CreateLocation(location *model.Location) error {
	_, err := r.collection.InsertOne(context.Background(), location)
	return err
}

// GetLocationByID retrieves a location by ID from the database.
func (r *LocationMongoRepository) GetLocationByID(id string) (*model.Location, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var location model.Location
	err = r.collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&location)
	if err != nil {
		return nil, err
	}
	return &location, nil
}

// UpdateLocation updates an existing location in the database.
func (r *LocationMongoRepository) UpdateLocation(id string, updatedLocation *model.Location) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updatedLocation})
	return err
}

// DeleteLocation removes a location from the database by ID.
func (r *LocationMongoRepository) DeleteLocation(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

// ListAll retrieves a list of all locations from the database.
func (r *LocationMongoRepository) ListAll() ([]model.Location, error) {
	var locations []model.Location
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &locations)
	if err != nil {
		return nil, err
	}
	return locations, nil
}