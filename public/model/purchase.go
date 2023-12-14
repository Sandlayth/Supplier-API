package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	ID         primitive.ObjectID    `json:"id" bson:"_id,omitempty"`
	Quantity   int                   `json:"quantity"`
	Date       time.Time             `json:"date"`
	Fees       float64               `json:"fees"`
	TotalPrice float64               `json:"totalPrice"`
	UserID     primitive.ObjectID    `json:"user" bson:"user"`
	LocationID   primitive.ObjectID  `json:"location" bson:"location"`
}