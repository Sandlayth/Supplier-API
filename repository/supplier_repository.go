package repository

import "github.com/sandlayth/supplier-api/model"

type SupplierRepository interface {
	CreateSupplier(supplier *model.Supplier) error
	GetSupplierByID(id string) (*model.Supplier, error)
	UpdateSupplier(id string, updatedSupplier *model.Supplier) error
	DeleteSupplier(id string) error
	ListAll() ([]model.Supplier, error)
}
