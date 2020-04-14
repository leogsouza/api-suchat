package handler

import (
	"log"
	"net/http"
	"time"

	gosocketio "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
	"github.com/leogsouza/api-suchat/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	UserImage string `json:"userImage"`
	NowTime   string `json:"nowTime"`
	Type      string `json:"type"`
	Message   string `json:"message"`
}

func (h *handler) socketHandler() {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected")
		//join them to room
		c.Join("chat")
	})

	//handle custom event
	server.On("input_message", func(c *gosocketio.Channel, msg *Message) string {
		log.Printf("%v", msg)
		userID, err := primitive.ObjectIDFromHex(msg.UserID)
		if err != nil {
			return err.Error()
		}

		createdAt, _ := time.Parse("2006-01-02T15:04:05Z07:00", msg.NowTime)

		chat := service.Chat{
			ID:        primitive.NewObjectID(),
			Message:   msg.Message,
			Sender:    userID,
			Type:      msg.Type,
			CreatedAt: createdAt,
		}

		out, err := h.SaveChat(chat)
		if err != nil {
			return err.Error()
		}
		//send event to all in room
		c.Emit("output_message", out)
		return "OK"
	})

	go func() {
		//setup http server
		serveMux := http.NewServeMux()
		serveMux.Handle("/socket.io/", server)
		log.Panic(http.ListenAndServe(":9999", serveMux))
	}()
}
