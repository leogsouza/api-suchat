package service

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	// ErrUserNotFound used when the user wasn't found on the db.
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidEmail used when the email is invalid.
	ErrInvalidEmail = errors.New("invalid email")
	// ErrInvalidUsername used when the username is invalid.
	ErrInvalidUsername = errors.New("invalid username")
	// ErrEmailTaken used when there is already an user registered with that email
	ErrEmailTaken = errors.New("email already exists")
	// ErrUsernameTaken used when there is already an user registered with that username
	ErrUsernameTaken = errors.New("username already exists")
	// ErrForbiddenFollow used when you try to follow yourself
	ErrForbiddenFollow = errors.New("cannot follow yourself")
	// ErrUnsupportedAvatarFormat used when the avatar file extension is invalid
	ErrUnsupportedAvatarFormat = errors.New("only png and jpeg allowed as avatar")
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

type UserChat struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Email     string             `bson:"email" json:"email"`
	Lastname  string             `bson:"lastname" json:"lastname"`
	AvatarURL *string            `bson:"avatar_url" json:"avatar_url"`
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

func (s *Service) findUserByEmail(email string) (User, error) {

	collection := s.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	u := User{}
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
	if err != nil {

		return u, err
	}

	return u, nil
}

func (s *Service) findUserChatById(id primitive.ObjectID) (UserChat, error) {
	collection := s.db.Collection("users")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	u := UserChat{}

	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
	if err != nil {
		return u, err
	}

	if u.AvatarURL == nil {
		avatarURL := "https://i.pravatar.cc/100?u=" + u.ID.Hex()
		u.AvatarURL = &avatarURL
	}

	return u, nil
}
