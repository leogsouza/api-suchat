package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Password  string             `bson:"password" json:"-"`
	Lastname  string             `bson:"lastname" json:"lastname"`
	AvatarURL *string            `bson:"avatar_url" json:"avatar_url"`
	Token     string             `bson:"token" json:"token"`
	TokenExp  time.Time          `bson:"token_exp" json:"expires_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at,omitempty"`
}

type userCreatedOutput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *Service) Register(name, email, lastname, password string) error {

	hashPassword, err := hash(password)
	if err != nil {
		return err
	}
	user := &User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		Lastname:  lastname,
		Password:  string(hashPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		TokenExp:  time.Time{},
	}

	collection := s.db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	collection.InsertOne(ctx, user)
	// send mail routine
	return nil
}

func hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
