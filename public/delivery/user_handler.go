package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sandlayth/supplier-api/public/entity"
	"github.com/sandlayth/supplier-api/public/usecase"
)

type UserHandler struct {
	useCase usecase.UserUseCase
}

func NewUserHandler(useCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{useCase: useCase}
}

func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody entity.User
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Continue with the registration logic...
	userID, err := h.useCase.Register(requestBody.Email, requestBody.Password, requestBody.FirstName, requestBody.LastName)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	response := map[string]string{"userID": userID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	token, err := h.useCase.Login(requestBody.Email, requestBody.Password)
	if err != nil {
		http.Error(w, "Failed to login", http.StatusUnauthorized)
		return
	}

	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request (@TODO later use authentication)
	userID := mux.Vars(r)["id"]

	err := h.useCase.Logout(userID)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request (@TODO later use authentication)
	userID := mux.Vars(r)["id"]

	userInfo, err := h.useCase.GetUserInfo(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func (h *UserHandler) ListAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from the request (@TODO later use authentication)
	users, err := h.useCase.ListAll()
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
