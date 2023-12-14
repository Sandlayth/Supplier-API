package route

import "github.com/gorilla/mux"

func AddUserRoutes(r *mux.Router, handler *UserHandler) {
	r.HandleFunc("/users/register", handler.RegisterHandler).Methods("POST")
	r.HandleFunc("/users/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/users/logout", handler.LogoutHandler).Methods("POST")
	r.HandleFunc("/users/{id}", handler.GetUserInfoHandler).Methods("GET")
	r.HandleFunc("/users", handler.ListAllUsersHandler).Methods("GET")
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
