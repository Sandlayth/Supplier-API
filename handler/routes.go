package handler

import (
	"github.com/gorilla/mux"
	"github.com/sandlayth/supplier-api/helper"
)

func AddUserRoutes(r *mux.Router, handler *UserHandler) {
	r.HandleFunc("/users/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/users/{id}/renew-token", handler.RenewTokenHandler).Methods("POST")
	// Admin Routes
	adminRouter := r.PathPrefix("/users").Subrouter()
	adminRouter.Use(helper.AdminAuthorizationMiddleware)
	adminRouter.HandleFunc("", handler.CreateUserHandler).Methods("POST")
	adminRouter.HandleFunc("/{id}", handler.UpdateUserHandler).Methods("PUT")
	adminRouter.HandleFunc("/{id}", handler.DeleteUserHandler).Methods("DELETE")
	adminRouter.HandleFunc("", handler.ListUsersHandler).Methods("GET")
	adminRouter.HandleFunc("/{id}", handler.GetUserHandler).Methods("GET")
	//	r.HandleFunc("/logout", handler.LogoutHandler).Methods("POST")
}

func AddLocationRoutes(r *mux.Router, handler *LocationHandler) {
	adminRouter := r.PathPrefix("/locations").Subrouter()
	adminRouter.Use(helper.AdminAuthorizationMiddleware)
	adminRouter.HandleFunc("/{id}", handler.UpdateLocationHandler).Methods("PUT")
	adminRouter.HandleFunc("/{id}", handler.DeleteLocationHandler).Methods("DELETE")
	adminRouter.HandleFunc("", handler.CreateLocationHandler).Methods("POST")

	managerRouter := r.PathPrefix("/locations").Subrouter()
	managerRouter.Use(helper.ManagerAuthorizationMiddleware)
	managerRouter.HandleFunc("/{id}", handler.GetLocationByIDHandler).Methods("GET")
	managerRouter.HandleFunc("", handler.ListAllLocationsHandler).Methods("GET")
	managerRouter.HandleFunc("/supplier/{id}", handler.ListBySupplierHandler).Methods("GET")
}

func AddSupplierRoutes(r *mux.Router, handler *SupplierHandler) {
	adminRouter := r.PathPrefix("/suppliers").Subrouter()
	adminRouter.Use(helper.AdminAuthorizationMiddleware)
	adminRouter.HandleFunc("", handler.CreateSupplierHandler).Methods("POST")
	adminRouter.HandleFunc("/{id}", handler.UpdateSupplierHandler).Methods("PUT")
	adminRouter.HandleFunc("/{id}", handler.DeleteSupplierHandler).Methods("DELETE")

	managerRouter := r.PathPrefix("/suppliers").Subrouter()
	managerRouter.Use(helper.ManagerAuthorizationMiddleware)
	managerRouter.HandleFunc("/{id}", handler.GetSupplierByIDHandler).Methods("GET")
	managerRouter.HandleFunc("", handler.GetAllSuppliersHandler).Methods("GET")
}

// AddPurchaseRoutes adds purchase-related routes to the provided router.
func AddPurchaseRoutes(r *mux.Router, handler *PurchaseHandler) {
	adminRouter := r.PathPrefix("/purchases").Subrouter()
	adminRouter.Use(helper.AdminAuthorizationMiddleware)
	adminRouter.HandleFunc("/{id}", handler.UpdatePurchaseHandler).Methods("PUT")
	adminRouter.HandleFunc("/{id}", handler.DeletePurchaseHandler).Methods("DELETE")
	adminRouter.HandleFunc("/user/{userID}", handler.ListPurchasesByUserHandler).Methods("GET")

	managerRouter := r.PathPrefix("/purchases").Subrouter()
	managerRouter.Use(helper.ManagerAuthorizationMiddleware)
	managerRouter.HandleFunc("", handler.CreatePurchaseHandler).Methods("POST")
	managerRouter.HandleFunc("/{id}", handler.GetPurchaseHandler).Methods("GET")
	managerRouter.HandleFunc("", handler.ListAllPurchasesHandler).Methods("GET")
}
