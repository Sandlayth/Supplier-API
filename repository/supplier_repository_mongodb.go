package repository

import (
	"context"

	"github.com/sandlayth/supplier-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SupplierMongoRepository struct {
	suppliersCollection *mongo.Collection
	locationsCollection *mongo.Collection
}

func NewSupplierMongoRepository(db *mongo.Database) *SupplierMongoRepository {
	return &SupplierMongoRepository{
		suppliersCollection: db.Collection("suppliers"),
		locationsCollection: db.Collection("locations"),
	}
}

// GetSupplierByID retrieves a supplier by ID from the database.
func (r *SupplierMongoRepository) GetSupplierByID(id string) (*model.Supplier, error) {
	var supplier model.Supplier
	idSupplier, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = r.suppliersCollection.FindOne(context.Background(), bson.M{"_id": idSupplier}).Decode(&supplier)
	if err != nil {
		return nil, err
	}
	return &supplier, nil
}

// GetAllSuppliers retrieves a list of all suppliers from the database.
func (r *SupplierMongoRepository) ListAll() ([]model.Supplier, error) {
	var suppliers []model.Supplier
	cursor, err := r.suppliersCollection.Find(context.Background(), bson.M{})
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
	_, err := r.suppliersCollection.InsertOne(context.Background(), supplier)
	return err
}

// UpdateSupplier updates an existing supplier in the database.
func (r *SupplierMongoRepository) UpdateSupplier(id string, updatedSupplier *model.Supplier) error {
	idSupplier, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.suppliersCollection.UpdateOne(context.Background(), bson.M{"_id": idSupplier}, bson.M{"$set": updatedSupplier})
	return err
}

// DeleteSupplier removes a supplier from the database by ID.
func (r *SupplierMongoRepository) DeleteSupplier(id string) error {
	idSupplier, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// First delete all locations related to the Supplier
	_, err = r.locationsCollection.DeleteMany(context.Background(), bson.M{"supplier": idSupplier})
	if err != nil {
		return err
	}
	// Then delete the Supplier
	_, err = r.suppliersCollection.DeleteOne(context.Background(), bson.M{"_id": idSupplier})
	return err
}
