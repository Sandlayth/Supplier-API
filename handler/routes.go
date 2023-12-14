package handler

import (
	"github.com/gorilla/mux"
)

func AddUserRoutes(r *mux.Router, handler *UserHandler) {
	r.HandleFunc("/users", handler.CreateUserHandler).Methods("POST")
	r.HandleFunc("/users/{id}", handler.GetUserHandler).Methods("GET")
	r.HandleFunc("/users/{id}", handler.UpdateUserHandler).Methods("PUT")
	r.HandleFunc("/users/{id}", handler.DeleteUserHandler).Methods("DELETE")
	r.HandleFunc("/users", handler.ListUsersHandler).Methods("GET")
//	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")
//	r.HandleFunc("/logout", handler.LogoutHandler).Methods("POST")
}

func AddLocationRoutes(r *mux.Router, handler *LocationHandler) {
	//r.HandleFunc("/locations/supplier/{supplierID}", handler.GetAllLocationsForSupplierHandler).Methods("GET")
	r.HandleFunc("/locations", handler.CreateLocationHandler).Methods("POST")
	r.HandleFunc("/locations/{id}", handler.GetLocationByIDHandler).Methods("GET")
	r.HandleFunc("/locations/{id}", handler.UpdateLocationHandler).Methods("PUT")
	r.HandleFunc("/locations/{id}", handler.DeleteLocationHandler).Methods("DELETE")
	r.HandleFunc("/locations", handler.ListAllLocationsHandler).Methods("GET")
	r.HandleFunc("/locations/supplier/{id}", handler.ListBySupplierHandler).Methods("GET")
}

func AddSupplierRoutes(r *mux.Router, handler *SupplierHandler) {
	r.HandleFunc("/suppliers/{id}", handler.GetSupplierByIDHandler).Methods("GET")
	r.HandleFunc("/suppliers", handler.GetAllSuppliersHandler).Methods("GET")
	r.HandleFunc("/suppliers", handler.CreateSupplierHandler).Methods("POST")
	r.HandleFunc("/suppliers/{id}", handler.UpdateSupplierHandler).Methods("PUT")
	r.HandleFunc("/suppliers/{id}", handler.DeleteSupplierHandler).Methods("DELETE")
}

// AddPurchaseRoutes adds purchase-related routes to the provided router.
func AddPurchaseRoutes(r *mux.Router, handler *PurchaseHandler) {
	r.HandleFunc("/purchases", handler.CreatePurchaseHandler).Methods("POST")
	r.HandleFunc("/purchases/{id}", handler.GetPurchaseHandler).Methods("GET")
	r.HandleFunc("/purchases/{id}", handler.UpdatePurchaseHandler).Methods("PUT")
	r.HandleFunc("/purchases/{id}", handler.DeletePurchaseHandler).Methods("DELETE")
	r.HandleFunc("/purchases", handler.ListAllPurchasesHandler).Methods("GET")
	r.HandleFunc("/purchases/user/{userID}", handler.ListPurchasesByUserHandler).Methods("GET")
}
