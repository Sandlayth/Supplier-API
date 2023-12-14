package model

import (
   "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
   ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
   Email     string             `json:"email"`
   Password  string             `json:"password"`
   FirstName string             `json:"first_name"`
   LastName  string             `json:"last_name"`
   Role      string             `json:"role"`
}