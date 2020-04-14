package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/leogsouza/api-suchat/internal/handler"
	"github.com/leogsouza/api-suchat/internal/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	port := env("PORT", "8080")

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	db := client.Database("suchat")
	s := service.New(db)

	h := handler.New(s)

	log.Printf("accepting connections on port %s", port)
	if err = http.ListenAndServe(":"+port, h); err != nil {
		log.Fatalf("could not start server: %v\n", err)
	}
}

func env(key, fallbackValue string) string {
	s := os.Getenv(key)
	if s == "" {
		return fallbackValue
	}

	return s
}
