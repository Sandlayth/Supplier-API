package entity

import (
   "go.mongodb.org/mongo-driver/bson/primitive"
)

type UserDAO struct {
   ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
   Email     string             `json:"email"`
   Password  string             `json:"password"`
   FirstName string             `json:"first_name"`
   LastName  string             `json:"last_name"`
}

// Convert User to UserDAO
func NewUserDAO(user *User) *UserDAO {
   return &UserDAO{
      ID:        primitive.NewObjectID(),
      Email:     user.Email,
      Password:  user.Password,
      FirstName: user.FirstName,
      LastName:  user.LastName,
   }
}

// Convert UserDAO to User
func (dao *UserDAO) ToUser() *User {
   return &User{
      ID:        dao.ID.Hex(),
      Email:     dao.Email,
      Password:  dao.Password,
      FirstName: dao.FirstName,
      LastName:  dao.LastName,
   }
}