package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

const (
	// TokenLifeSpan until tokens are valid
	TokenLifeSpan = time.Hour * 24 * 14
	// KeyAuthUserID to use in context
	KeyAuthUserID key = "auth_user_id"
)

var (
	// ErrUnauthenticated used when there is no user authenticated in the context.
	ErrUnauthenticated = errors.New("unauthenticated")
)

type key string

// LoginOutput response
type UserLoginOutput struct {
	Name      string    `bson:"name" json:"name"`
	Email     string    `bson:"email" json:"email"`
	Lastname  string    `bson:"lastname" json:"lastname"`
	AvatarURL *string   `bson:"avatar_url" json:"avatar_url"`
	Token     string    `bson:"token" json:"token"`
	TokenExp  time.Time `bson:"token_exp" json:"expires_at"`
}

func (s *Service) Login(email, password string) (UserLoginOutput, error) {

	var out UserLoginOutput

	collection := s.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	u := User{}
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {

		return out, err
	}

	if err := verifyPassword(u.Password, password); err != nil {
		return out, err
	}
	out.Name = u.Name
	out.Email = u.Email
	out.Lastname = u.Lastname
	out.AvatarURL = u.AvatarURL

	err = s.generateToken(u.Email, &out)

	return out, nil
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *Service) generateToken(email string, u *UserLoginOutput) error {
	tokenExp := time.Now().Add(time.Hour * 1)
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["email"] = email
	claims["exp"] = tokenExp //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	if err != nil {
		return nil
	}

	// Update user token and tokenExp
	collection := s.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	update := bson.D{{"$set", bson.M{"token": tokenStr, "token_exp": tokenExp}}}
	err = collection.FindOneAndUpdate(ctx, bson.M{"email": email}, update).Decode(&u)
	if err != nil {
		return err
	}

	u.Token = tokenStr
	u.TokenExp = tokenExp

	return nil

}
