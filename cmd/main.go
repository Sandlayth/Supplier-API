package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sandlayth/supplier-api/handler"
	"github.com/sandlayth/supplier-api/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client := initDb()
	defer client.Disconnect(context.Background())

	// Select the database and initialize the repositories
	db := client.Database("supplier-api")
	userRepo := repository.NewUserMongoRepository(db)
	locationRepo := repository.NewLocationMongoRepository(db)
	supplierRepo := repository.NewSupplierMongoRepository(db)
	purchaseRepo := repository.NewPurchaseMongoRepository(db)

	// Initialize the handlers
	userHandler := handler.NewUserHandler(userRepo)
	locationHandler := handler.NewLocationHandler(locationRepo)
	supplierHandler := handler.NewSupplierHandler(supplierRepo)
	purchaseHandler := handler.NewPurchaseHandler(purchaseRepo)

	// Initialize the router and add the routes
	router := mux.NewRouter()
	handler.AddUserRoutes(router, userHandler)
	handler.AddLocationRoutes(router, locationHandler)
	handler.AddSupplierRoutes(router, supplierHandler)
	handler.AddPurchaseRoutes(router, purchaseHandler)

	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		})

	corsRouter := cors.Handler(router)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", corsRouter))
}

func initDb() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	client.Database("supplier-api")

	return client
}
