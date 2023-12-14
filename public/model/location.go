package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Location struct {
	ID       primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	SupplierID primitive.ObjectID  `json:"supplier" bson:"supplier"`
}
