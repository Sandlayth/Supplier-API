package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sandlayth/supplier-api/helper"
	"github.com/sandlayth/supplier-api/model"
	"github.com/sandlayth/supplier-api/repository"
)

// PurchaseHandler handles HTTP requests related to purchases.
type PurchaseHandler struct {
	pr repository.PurchaseRepository
}

// NewPurchaseHandler creates a new instance of PurchaseHandler.
func NewPurchaseHandler(pr repository.PurchaseRepository) *PurchaseHandler {
	return &PurchaseHandler{pr: pr}
}

// CreatePurchaseHandler handles requests to create a new purchase.
func (h *PurchaseHandler) CreatePurchaseHandler(w http.ResponseWriter, r *http.Request) {
	var purchase model.Purchase
	err := json.NewDecoder(r.Body).Decode(&purchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	claims, ok := r.Context().Value("userClaims").(*model.Claims)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	purchase.UserID = claims.UserID
	err = h.pr.CreatePurchase(&purchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, purchase)
}

// GetPurchaseHandler handles requests to retrieve a purchase by ID.
func (h *PurchaseHandler) GetPurchaseHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	purchaseID := params["id"]

	purchase, err := h.pr.GetPurchaseByID(purchaseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if purchase == nil {
		http.NotFound(w, r)
		return
	}
	claims, ok := r.Context().Value("userClaims").(*model.Claims)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if claims.UserID.Hex() != purchase.ID.Hex() {
		helper.RespondJSON(w, nil)
	}
	helper.RespondJSON(w, purchase)
}

// UpdatePurchaseHandler handles requests to update a purchase by ID.
func (h *PurchaseHandler) UpdatePurchaseHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	purchaseID := params["id"]

	var updatedPurchase model.Purchase
	err := json.NewDecoder(r.Body).Decode(&updatedPurchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.pr.UpdatePurchase(purchaseID, &updatedPurchase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, updatedPurchase)
}

// DeletePurchaseHandler handles requests to delete a purchase by ID.
func (h *PurchaseHandler) DeletePurchaseHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	purchaseID := params["id"]

	err := h.pr.DeletePurchase(purchaseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, map[string]string{"message": "Purchase deleted successfully"})
}

// ListAllPurchasesHandler handles requests to retrieve a list of all purchases.
func (h *PurchaseHandler) ListAllPurchasesHandler(w http.ResponseWriter, r *http.Request) {
	var purchases []model.Purchase
	var err error
	claims, ok := r.Context().Value("userClaims").(*model.Claims)
	if !ok {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if claims.Role == "admin" {
		purchases, err = h.pr.ListAll()
	} else {
		purchases, err = h.pr.ListPurchasesByUser(claims.UserID.Hex())
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, purchases)
}

// ListPurchasesByUserHandler handles requests to retrieve a list of purchases for a specific user.
func (h *PurchaseHandler) ListPurchasesByUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["userID"]

	purchases, err := h.pr.ListPurchasesByUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, purchases)
}
