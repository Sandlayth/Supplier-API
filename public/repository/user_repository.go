package repository

import "github.com/sandlayth/supplier-api/public/entity"

type UserRepository interface {
	Create(user *entity.User) (string, error)
	FindByEmail(email string) (*entity.User, error)
	FindById(id string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id string) error
	ListAll() (*[]entity.User, error)
}
