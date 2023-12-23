package repository

import "github.com/sandlayth/supplier-api/model"

type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(id string, updatedUser *model.User) error
	DeleteUser(id string) error
	ListAll() (*[]model.User, error)
	GetTokens(user *model.User) (string, string, error)
	RenewTokens(userID string, refreshToken string) (string, string, error)
    ValidateUserCredentials(user *model.User) error
	//RevokeToken(userID string, refreshToken string) error
}
