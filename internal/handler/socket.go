package handler

import (
	"github.com/ambelovsky/gosf"
)

func init() {
	// Listen on an endpoint
	gosf.Listen("chat", chat)
}

type ChatRequestBody struct {
	UserID    string `json:"userId"`
	Username  string `json:"username"`
	UserImage string `json:"userImage"`
	NowTime   string `json:"nowTime"`
	Type      string `json:"type"`
}

func chat(client *gosf.Client, request *gosf.Request) *gosf.Message {

	requestBody := new(ChatRequestBody)
	gosf.MapToStruct(request.Message.Body, requestBody)

	return gosf.NewSuccessMessage(request.Message.Text, gosf.StructToMap(requestBody))
}

func SocketHandler() {
	go func() {
		// Start the server using a basic configuration
		gosf.Startup(map[string]interface{}{"port": 9999})
	}()
}
