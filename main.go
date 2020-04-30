package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/leogsouza/api-suchat/internal/handler"
	"github.com/leogsouza/api-suchat/internal/helper"
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

	port := helper.Env("PORT", "8080")
	dbHost := helper.Env("DB_HOST", "mongodb://localhost")
	dbPort := helper.Env("DB_PORT", "27017")

	clientOptions := options.Client().ApplyURI(fmt.Sprintf("%s:%s", dbHost, dbPort))

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
