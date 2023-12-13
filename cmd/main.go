package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sandlayth/supplier-api/public/delivery"
	"github.com/sandlayth/supplier-api/public/repository"
	"github.com/sandlayth/supplier-api/public/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	// Select the database and initialize the repository
	db := client.Database("supplier-api")
	repo := repository.NewMongoRepository(db)

	// Initialize the use case
	useCase := usecase.NewUserUseCase(repo)

	// Initialize the handler
	handler := delivery.NewUserHandler(useCase)

	// Initialize the router
	router := delivery.NewRouter(handler)

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8080", router))
}
