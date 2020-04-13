package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// ErrInvalidToken used when there is no user authenticated in the context.
	ErrInvalidToken = errors.New("InvalidToken")
)

type key string

// LoginOutput response
type UserLoginOutput struct {
	ID           primitive.ObjectID `bson:"_id" json:"userId,omitempty"`
	Name         string             `bson:"name" json:"name"`
	Email        string             `bson:"email" json:"email"`
	Lastname     string             `bson:"lastname" json:"lastname"`
	AvatarURL    *string            `bson:"avatar_url" json:"avatar_url"`
	Token        string             `bson:"token" json:"token"`
	TokenExp     time.Time          `bson:"token_exp" json:"expires_at"`
	LoginSuccess bool               `json:"loginSuccess"`
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
	out.ID = u.ID
	out.Name = u.Name
	out.Email = u.Email
	out.Lastname = u.Lastname
	out.AvatarURL = u.AvatarURL
	out.LoginSuccess = true

	err = s.generateToken(u.Email, &out)

	return out, nil
}

func (s *Service) Logout(ctx context.Context) error {
	email, ok := ctx.Value(KeyAuthUserID).(string)
	if !ok {
		return ErrUnauthenticated
	}
	collection := s.db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{"email": email}
	update := bson.D{{"$set", bson.M{"token": "", "token_exp": time.Time{}}}}
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
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
		return err
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

// AuthUserEmailID retrieves the user ID from the token
func (s *Service) AuthUserEmailID(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		email := claims["email"].(string)
		expiresAt := claims["exp"].(string)
		timeExp, _ := time.Parse("2006-01-02T15:04:05-07:00", expiresAt)
		if err = s.validateUserTokenExpiration(tokenString, email, timeExp); err != nil {
			return "", err
		}
		return email, nil
	}
	return "", nil

}

type authResponse struct {
	User
	IsAuth bool `json:"isAuth"`
	Error  bool `json:"error"`
}

// AuthUser retrieves user from the context
func (s *Service) AuthUser(ctx context.Context) (authResponse, error) {

	var resp authResponse
	uid, ok := ctx.Value(KeyAuthUserID).(string)
	log.Println("uid ok", uid, ok)
	if !ok {
		return resp, ErrUnauthenticated
	}
	u, err := s.findUserByEmail(uid)
	if err != nil {
		return resp, err
	}
	resp.User = u
	resp.IsAuth = true
	resp.Error = false
	return resp, nil
}

func (s *Service) validateUserTokenExpiration(token, email string, expiresAt time.Time) error {

	now := time.Now()
	if now.After(expiresAt) {
		return ErrInvalidToken
	}
	collection := s.db.Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	u := User{}
	err := collection.FindOne(ctx, bson.M{"token": token, "email": email, "token_exp": expiresAt}).Decode(&u)
	if err != nil {
		return err
	}

	return nil
}
