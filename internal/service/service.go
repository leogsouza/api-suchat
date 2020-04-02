package service

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db *mongo.Database
}

func New(database *mongo.Database) *Service {

	return &Service{
		db: database,
	}
}
