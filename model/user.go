package model

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
   ID            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
   Email         string             `json:"email"`
   Password      string             `json:"password"`
   FirstName     string             `json:"first_name"`
   LastName      string             `json:"last_name"`
   Role          string             `json:"role"`
   *Claims
}

type Claims struct {
   UserID   primitive.ObjectID      `json:"userID"`
   Role string                      `json:"role" bson:"omitempty"`
   jwt.RegisteredClaims
}