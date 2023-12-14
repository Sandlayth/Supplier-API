package repository

import "github.com/sandlayth/supplier-api/model"

type PurchaseRepository interface {
	CreatePurchase(purchase *model.Purchase) error
	GetPurchaseByID(id string) (*model.Purchase, error)
	UpdatePurchase(id string, updatedPurchase *model.Purchase) error
	DeletePurchase(id string) error
	ListAll() ([]model.Purchase, error)
	ListPurchasesByUser(user string) ([]model.Purchase, error)
}