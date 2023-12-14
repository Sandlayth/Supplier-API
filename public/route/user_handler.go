package route

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sandlayth/supplier-api/public/model"
	"github.com/sandlayth/supplier-api/public/repository"
)

type UserHandler struct {
	ur repository.UserRepository
}

func NewUserHandler(r repository.UserRepository) *UserHandler {
	return &UserHandler{ur: r}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody model.User
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Continue with the registration logic...
	userID, err := h.ur.Register(requestBody.Email, requestBody.Password, requestBody.FirstName, requestBody.LastName)
	if err != nil {
		http.Error(w, "Failed to register user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	respondJSON(w, map[string]string{"userID": userID})
}

func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.ur.Login(requestBody.Email, requestBody.Password)
	if err != nil {
		http.Error(w, "Failed to login", http.StatusUnauthorized)
		return
	}

	respondJSON(w, map[string]string{"token": token})
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request (@TODO later use authentication)
	userID := mux.Vars(r)["id"]

	err := h.ur.Logout(userID)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request (@TODO later use authentication)
	userID := mux.Vars(r)["id"]

	userInfo, err := h.ur.GetUserInfo(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	respondJSON(w, userInfo)
}

func (h *UserHandler) ListAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request (@TODO later use authentication)
	users, err := h.ur.ListAll()
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	respondJSON(w, users)
}
