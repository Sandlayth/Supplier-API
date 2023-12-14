package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Supplier struct {
	ID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}