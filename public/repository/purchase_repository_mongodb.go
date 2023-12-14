package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/sandlayth/supplier-api/public/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PurchaseMongoRepository is a concrete implementation of PurchaseRepository using MongoDB.
type PurchaseMongoRepository struct {
	purchasesCollection *mongo.Collection
	locationsCollection *mongo.Collection
	usersCollection     *mongo.Collection
}

func NewPurchaseMongoRepository(db *mongo.Database) *PurchaseMongoRepository {
	return &PurchaseMongoRepository{
		purchasesCollection: db.Collection("purchases"),
		locationsCollection: db.Collection("locations"),
		usersCollection: db.Collection("users"),
	}
}

// CreatePurchase adds a new purchase to the database.
func (r *PurchaseMongoRepository) CreatePurchase(purchase *model.Purchase) error {
	// Validate that the specified UserID corresponds to an existing user
	if err := r.validateUser(purchase.UserID); err != nil {
		return err
	}
	
	totalPrice, err := r.calculatePrice(purchase)
	if err != nil {
		return err
	} 
	purchase.TotalPrice = totalPrice
	
	// Continue with purchase creation
	result, err := r.purchasesCollection.InsertOne(context.Background(), purchase)

	// Update the purchase with the new ID
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Printf("Error converting inserted ID to primitive.ObjectID\n")
		return errors.New("inserted ID is not a primitive.ObjectID")
	}
	purchase.ID = insertedID
	return err
}

// GetPurchaseByID retrieves a purchase by ID from the database.
func (r *PurchaseMongoRepository) GetPurchaseByID(id string) (*model.Purchase, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var purchase model.Purchase
	err = r.purchasesCollection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&purchase)
	if err != nil {
		return nil, err
	}
	return &purchase, nil
}

// UpdatePurchase updates an existing purchase in the database.
func (r *PurchaseMongoRepository) UpdatePurchase(id string, updatedPurchase *model.Purchase) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if err := r.validateUser(updatedPurchase.UserID); err != nil {
		return err
	}
	
	totalPrice, err := r.calculatePrice(updatedPurchase)
	if err != nil {
		return err
	} 
	updatedPurchase.TotalPrice = totalPrice

	_, err = r.purchasesCollection.UpdateOne(context.Background(), bson.M{"_id": objectID}, bson.M{"$set": updatedPurchase})
	return err
}

// DeletePurchase removes a purchase from the database by ID.
func (r *PurchaseMongoRepository) DeletePurchase(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.purchasesCollection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

// ListAll retrieves a list of all purchases from the database.
func (r *PurchaseMongoRepository) ListAll() ([]model.Purchase, error) {
	var purchases []model.Purchase
	cursor, err := r.purchasesCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &purchases)
	if err != nil {
		return nil, err
	}
	return purchases, nil
}

// ListPurchasesByUser retrieves a list of purchases for a specific user from the database.
func (r *PurchaseMongoRepository) ListPurchasesByUser(user string) ([]model.Purchase, error) {
	userID, err := primitive.ObjectIDFromHex(user)
	if err != nil {
		return nil, err
	}
	if err := r.validateUser(userID); err != nil {
		return nil, err
	}
	
	var purchases []model.Purchase
	cursor, err := r.purchasesCollection.Find(context.Background(), bson.M{"user": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	err = cursor.All(context.Background(), &purchases)
	if err != nil {
		return nil, err
	}
	return purchases, nil
}

// calculatePrice calculate the price of the purchase (quantity * price * (1 - fees))
func (r *PurchaseMongoRepository) calculatePrice(purchase *model.Purchase) (float64, error) {
	// Retrieve the corresponding location to get the price
	location, err := r.getLocationByID(purchase.LocationID)
	if err != nil {
		return 0.0, err
	}
	price := float64(purchase.Quantity) * location.Price * (1 - purchase.Fees)
	return price, nil
}

// validateUser checks if a user with the given ID exists.
func (r *PurchaseMongoRepository) validateUser(userID primitive.ObjectID) error {
	// Check if the specified UserID exists in the User collection
	userCount, err := r.usersCollection.CountDocuments(context.Background(), bson.M{"_id": userID})
	if err != nil {
		return err
	}
	if userCount == 0 {
		return fmt.Errorf("user with ID %s does not exist", userID.Hex())
	}
	return nil
}

// getLocationByID retrieves a location by ID from the database.
func (r *PurchaseMongoRepository) getLocationByID(locationID primitive.ObjectID) (*model.Location, error) {
	var location model.Location
	err := r.locationsCollection.FindOne(context.Background(), bson.M{"_id": locationID}).Decode(&location)
	if err != nil {
		return nil, err
	}
	return &location, nil
}