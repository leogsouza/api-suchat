package service

import (
	"context"
	"time"

	"github.com/leogsouza/api-suchat/internal/model"
)

type userCreatedOutput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *Service) Register(user model.User) {
	collection := s.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	collection.InsertOne(ctx, user)
}
