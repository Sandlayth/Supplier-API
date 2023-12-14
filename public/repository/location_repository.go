package repository

import "github.com/sandlayth/supplier-api/public/model"

type LocationRepository interface {
	//GetAllLocationsForLocation(supplierID string) ([]model.Location, error)
	CreateLocation(supplier *model.Location) error
	GetLocationByID(id string) (*model.Location, error)
	UpdateLocation(id string, updatedLocation *model.Location) error
	DeleteLocation(id string) error
	ListAll() ([]model.Location, error)
}