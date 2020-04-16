package service

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Chat struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Sender    primitive.ObjectID `bson:"sender" json:"sender"`
	Message   string             `bson:"message" json:"message"`
	Type      string             `bson:"type" json:"type"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at,omitempty"`
}

type chatOutput struct {
	ID        primitive.ObjectID `json:"id"`
	Sender    UserChat           `json:"sender"`
	Message   string             `json:"message"`
	Type      string             `json:"type"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
}

func (s *Service) SaveChat(c Chat) (Chat, error) {
	collection := s.db.Collection("chats")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result, err := collection.InsertOne(ctx, c)
	var chat Chat
	if err != nil {
		return chat, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&chat)
	if err != nil {
		return chat, err
	}
	return chat, nil
}

func (s *Service) GetChats() ([]chatOutput, error) {

	var chats []chatOutput
	collection := s.db.Collection("chats")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return chats, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var chat Chat
		err = cur.Decode(&chat)
		if err != nil {
			log.Fatal(err)
		}

		u, err := s.findUserChatById(chat.Sender)
		if err != nil {
			return chats, err
		}

		chout := chatOutput{
			chat.ID,
			u,
			chat.Message,
			chat.Type,
			chat.CreatedAt,
		}

		chats = append(chats, chout)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return chats, nil
}
