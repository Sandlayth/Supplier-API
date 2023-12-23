package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sandlayth/supplier-api/helper"
	"github.com/sandlayth/supplier-api/model"
	"github.com/sandlayth/supplier-api/repository"
)

type UserHandler struct {
	ur repository.UserRepository
}

func NewUserHandler(r repository.UserRepository) *UserHandler {
	return &UserHandler{ur: r}
}

// CreateUserHandler handles requests to create a new user.
func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.ur.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, user)
}

// GetUserHandler handles requests to retrieve a user by ID.
func (h *UserHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	user, err := h.ur.GetUserByID(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.NotFound(w, r)
		return
	}

	helper.RespondJSON(w, user)
}

// UpdateUserHandler handles requests to update a user by ID.
func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	var updatedUser model.User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.ur.UpdateUser(userID, &updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, updatedUser)
}

// DeleteUserHandler handles requests to delete a user by ID.
func (h *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userID := params["id"]

	err := h.ur.DeleteUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, map[string]string{"message": "User deleted successfully"})
}

// ListUsersHandler handles requests to retrieve a list of all users.
func (h *UserHandler) ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := h.ur.ListAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, users)
}

// LoginHandler handles requests for user login and generates an authentication token.
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.ur.ValidateUserCredentials(&user); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	dbUser, err := h.ur.GetUserByEmail(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate authentication token
	refreshToken, accessToken, err := h.ur.GetTokens(dbUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	helper.RespondJSON(w, map[string]string{"refresh_token": refreshToken, "access_token": accessToken})
}

func (h *UserHandler) RenewTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	refreshToken := r.Header.Get("Authorization")

	// Renew tokens for the user
	newAccessToken, newRefreshToken, err := h.ur.RenewTokens(userID, refreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the new tokens
	helper.RespondJSON(w, map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

/*
// LogoutHandler handles requests to revoke the authentication token for a user.
func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	if userID == "" {
		http.Error(w, "UserID not provided in the request header", http.StatusBadRequest)
		return
	}

	err := h.ur.RevokeAuthToken(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helper.RespondJSON(w, map[string]string{"message": "Logout successful"})
}
*/
