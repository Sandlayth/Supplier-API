package repository

import (
	"context"

	"github.com/sandlayth/supplier-api/public/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SupplierMongoRepository struct {
	collection *mongo.Collection
}

func NewSupplierMongoRepository(db *mongo.Database) *SupplierMongoRepository {
	return &SupplierMongoRepository{
		collection: db.Collection("suppliers"),
	}
}

// GetSupplierByID retrieves a supplier by ID from the database.
func (r *SupplierMongoRepository) GetSupplierByID(id string) (*model.Supplier, error) {
	var supplier model.Supplier
	idSupplier, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.collection.FindOne(context.Background(), bson.M{"_id": idSupplier}).Decode(&supplier)
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

// GetAllSuppliers retrieves a list of all suppliers from the database.
func (r *SupplierMongoRepository) ListAll() ([]model.Supplier, error) {
	var suppliers []model.Supplier
	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &suppliers)
	if err != nil {
		return nil, err
	}
	return suppliers, nil
}

// CreateSupplier adds a new supplier to the database.
func (r *SupplierMongoRepository) CreateSupplier(supplier *model.Supplier) error {
	_, err := r.collection.InsertOne(context.Background(), supplier)
	return err
}

// UpdateSupplier updates an existing supplier in the database.
func (r *SupplierMongoRepository) UpdateSupplier(id string, updatedSupplier *model.Supplier) error {
	idSupplier, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.UpdateOne(context.Background(), bson.M{"_id": idSupplier}, bson.M{"$set": updatedSupplier})
	return err
}

// DeleteSupplier removes a supplier from the database by ID.
func (r *SupplierMongoRepository) DeleteSupplier(id string) error {
	idSupplier, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": idSupplier})
	return err
}
