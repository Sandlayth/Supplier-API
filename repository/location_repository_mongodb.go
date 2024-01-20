package repository

import (
	"context"
	"fmt"

	"github.com/sandlayth/supplier-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationMongoRepository struct {
	locationsCollection *mongo.Collection
	suppliersCollection *mongo.Collection
}

func NewLocationMongoRepository(db *mongo.Database) *LocationMongoRepository {
	return &LocationMongoRepository{
		locationsCollection: db.Collection("locations"),
		suppliersCollection: db.Collection("suppliers"),
	}
}

// CreateLocation adds a new location to the database.
func (r *LocationMongoRepository) CreateLocation(location *model.Location) error {
	err := r.supplierExists(location.SupplierID)
	if err != nil {
		return err
	}
	_, err = r.locationsCollection.InsertOne(context.Background(), location)
	return err
}

// GetLocationByID retrieves a location by ID from the database.
func (r *LocationMongoRepository) GetLocationByID(id string) (*model.Location, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var location model.Location
	err = r.locationsCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&location)
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
	err = r.supplierExists(updatedLocation.SupplierID)
	if err != nil {
		return err
	}

	_, err = r.locationsCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updatedLocation})
	return err
}

// DeleteLocation removes a location from the database by ID.
func (r *LocationMongoRepository) DeleteLocation(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.locationsCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

// ListAll retrieves a list of all locations from the database.
func (r *LocationMongoRepository) ListAll() ([]model.Location, error) {
	var locations []model.Location

	pipeline := bson.A{
		bson.D{{"$lookup", bson.D{{"from", "suppliers"}, {"localField", "supplier"}, {"foreignField", "_id"}, {"as", "supplierInfo"}}}},
		bson.D{{"$unwind", "$supplierInfo"}},
		bson.D{{"$project", bson.D{
			{"_id", 1},
			{"name", 1},
			{"price", 1},
			{"supplier", 1},
			{"supplierName", "$supplierInfo.name"},
		}}},
	}

	cursor, err := r.locationsCollection.Aggregate(context.Background(), pipeline)
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

// ListBySupplier retrieves a list of all locations for a specific supplier from the database.
func (r *LocationMongoRepository) ListBySupplier(id string) ([]model.Location, error) {
	var locations []model.Location

	supplierID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	cursor, err := r.locationsCollection.Find(context.Background(), bson.M{"supplier": supplierID})
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

func (r *LocationMongoRepository) supplierExists(supplierID primitive.ObjectID) (error) {
	count, err := r.suppliersCollection.CountDocuments(context.Background(), bson.M{"_id": supplierID})
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("supplier with ID %s does not exist", supplierID.Hex())
	}
	return nil
}
