package repository

import "github.com/sandlayth/supplier-api/model"

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(id string, updatedUser *model.User) error
	DeleteUser(id string) error
	ListAll() (*[]model.User, error)
	//GenerateAuthToken(user *model.User) (string, error)
	//RevokeAuthToken(userID string) error
}
