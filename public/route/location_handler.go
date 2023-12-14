package route

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sandlayth/supplier-api/public/model"
	"github.com/sandlayth/supplier-api/public/repository"
)

type LocationHandler struct {
	lr repository.LocationRepository
}

func NewLocationHandler(r repository.LocationRepository) *LocationHandler {
	return &LocationHandler{lr: r}
}

// CreateLocationHandler handles requests to create a new location.
func (h *LocationHandler) CreateLocationHandler(w http.ResponseWriter, r *http.Request) {
	var newLocation model.Location
	err := json.NewDecoder(r.Body).Decode(&newLocation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.lr.CreateLocation(&newLocation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "Location created successfully"})
}

// GetLocationByIDHandler handles requests to retrieve a location by ID.
func (h *LocationHandler) GetLocationByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	locationID := params["id"]

	location, err := h.lr.GetLocationByID(locationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, location)
}

// UpdateLocationHandler handles requests to update an existing location.
func (h *LocationHandler) UpdateLocationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	locationID := params["id"]

	var updatedLocation model.Location
	err := json.NewDecoder(r.Body).Decode(&updatedLocation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.lr.UpdateLocation(locationID, &updatedLocation)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "Location updated successfully"})
}

// DeleteLocationHandler handles requests to delete a location by ID.
func (h *LocationHandler) DeleteLocationHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	locationID := params["id"]

	err := h.lr.DeleteLocation(locationID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "Location deleted successfully"})
}

// GetAllLocationsHandler handles requests to retrieve all unique locations.
func (h *LocationHandler) ListAllLocationsHandler(w http.ResponseWriter, r *http.Request) {
	locations, err := h.lr.ListAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, locations)
}

// ListBySupplierHandler handles requests to retrieve a list of all locations for a specific supplier.
func (h *LocationHandler) ListBySupplierHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	supplierID := params["id"]

	locations, err := h.lr.ListBySupplier(supplierID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, locations)
}
