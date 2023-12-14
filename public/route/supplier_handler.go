package route

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sandlayth/supplier-api/public/model"
	"github.com/sandlayth/supplier-api/public/repository"
)

type SupplierHandler struct {
	sr repository.SupplierRepository
}

// NewSupplierHandler creates a new instance of SupplierHandler.
func NewSupplierHandler(r repository.SupplierRepository) *SupplierHandler {
	return &SupplierHandler{sr: r}
}

// GetSupplierByIDHandler handles requests to retrieve a supplier by ID.
func (h *SupplierHandler) GetSupplierByIDHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	supplierID := params["id"]

	supplier, err := h.sr.GetSupplierByID(supplierID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, supplier)
}

// GetAllSuppliersHandler handles requests to retrieve all suppliers.
func (h *SupplierHandler) GetAllSuppliersHandler(w http.ResponseWriter, r *http.Request) {
	suppliers, err := h.sr.ListAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, suppliers)
}

// CreateSupplierHandler handles requests to create a new supplier.
func (h *SupplierHandler) CreateSupplierHandler(w http.ResponseWriter, r *http.Request) {
	var newSupplier model.Supplier
	err := json.NewDecoder(r.Body).Decode(&newSupplier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.sr.CreateSupplier(&newSupplier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "Supplier created successfully"})
}

// UpdateSupplierHandler handles requests to update an existing supplier.
func (h *SupplierHandler) UpdateSupplierHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	supplierID := params["id"]

	var updatedSupplier model.Supplier
	err := json.NewDecoder(r.Body).Decode(&updatedSupplier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.sr.UpdateSupplier(supplierID, &updatedSupplier)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "Supplier updated successfully"})
}

// DeleteSupplierHandler handles requests to delete a supplier by ID.
func (h *SupplierHandler) DeleteSupplierHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	supplierID := params["id"]

	err := h.sr.DeleteSupplier(supplierID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"message": "Supplier deleted successfully"})
}