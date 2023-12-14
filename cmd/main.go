package main

import (
	"context"
	"log"
	"net/http"

	"github.com/sandlayth/supplier-api/public/route"
	"github.com/sandlayth/supplier-api/public/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func main() {
	client := initDb()
	defer client.Disconnect(context.Background())

	// Select the database and initialize the repository
	db := client.Database("supplier-api")
	userRepo := repository.NewMongoRepository(db)

	// Initialize the handler
	handler := route.NewUserHandler(userRepo)

	// Initialize the router
	router := route.NewRouter(handler)

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