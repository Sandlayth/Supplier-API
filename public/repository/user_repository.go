package repository

import "github.com/sandlayth/supplier-api/public/model"

type UserRepository interface {
	Register(email, password, firstName, lastName string) (string, error)
	Login(email, password string) (string, error)
	Logout(userID string) error
	GetUserInfo(userID string) (*model.User, error)
	ListAll() (*[]model.User, error)
}
