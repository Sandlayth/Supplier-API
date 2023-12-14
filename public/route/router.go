package route

import "github.com/gorilla/mux"

func NewRouter(handler *UserHandler) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/users/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/users/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/users/logout", handler.LogoutHandler).Methods("POST")
	r.HandleFunc("/users/{id}", handler.GetUserInfoHandler).Methods("GET")
	r.HandleFunc("/users", handler.ListAllUsersHandler).Methods("GET")

	return r
}