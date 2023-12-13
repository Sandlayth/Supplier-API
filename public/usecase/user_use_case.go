package usecase

import "github.com/sandlayth/supplier-api/public/entity"

type UserUseCase interface {
	Register(email, password, firstName, lastName string) (string, error)
	Login(email, password string) (string, error)
	Logout(userID string) error
	GetUserInfo(userID string) (*entity.User, error)
	ListAll() (*[]entity.User, error)
}
