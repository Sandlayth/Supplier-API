package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sandlayth/supplier-api/public/repository"
	"github.com/sandlayth/supplier-api/public/route"

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

	// Initialize the handlers
	userHandler := route.NewUserHandler(userRepo)
	locationHandler := route.NewLocationHandler(locationRepo)
	supplierHandler := route.NewSupplierHandler(supplierRepo)

	// Initialize the router and add the routes
	router := mux.NewRouter()
	route.AddUserRoutes(router, userHandler)
	route.AddLocationRoutes(router, locationHandler)
	route.AddSupplierRoutes(router, supplierHandler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", router))
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