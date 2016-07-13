package server

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"github.com/poslegm/blockchain-chat/network"
	"github.com/poslegm/blockchain-chat/db"
)

type WebSocketMessage struct {
	Type     string
	Messages []ChatMessage
}

type ChatMessage struct {
	Receiver string
	Sender string
	Text string
}
// TODO вкл/выкл майнинга
func receive(ws *websocket.Conn) {
	// чтение не должно прекращаться
	ws.SetReadDeadline(time.Time{})

	defer ws.Close()
	for {
		msg := WebSocketMessage{}

		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("WebSockets.receive json: " + err.Error())
			return
		}

		msg.switchTypes(ws)
	}
}

func (msg WebSocketMessage) switchTypes(ws *websocket.Conn) {
	switch msg.Type {
	case "GetMessages":
		networkMessages, err := db.GetAllMessages()
		if err != nil {
			fmt.Println("Websockets.switchTypes: ", err.Error())
			return
		}

		chatMessages := make([]ChatMessage, 0)
		for _, networkMsg := range networkMessages {
			chatMessages = append(chatMessages, ChatMessage{
				networkMsg.Receiver, networkMsg.Sender, networkMsg.Text,
			})
		}

		sendMessage(WebSocketMessage{"AllMessages", chatMessages}, ws)
	case "SendMessage":
		if len(msg.Messages) != 1 {
			fmt.Printf("WebSocket.switchTypes: incorrect message - %#v\n", msg)
			return
		}
		chatMsg := msg.Messages[0]
		network.CurrentNetworkUser.SendMessage(network.NetworkMessage{
			chatMsg.Receiver,
			chatMsg.Sender,
			chatMsg.Text,
		})
	}
}

func sendMessage(msg WebSocketMessage, ws *websocket.Conn) {
	go func() {
		err := ws.WriteJSON(&msg)
		if err != nil {
			fmt.Println("WebSocket.sendMessage: " + err.Error())
		}
	}()
}

func createConnection(ws *websocket.Conn) {
	go receive(ws)
}

func createWSHandler() http.HandlerFunc {
	return func (w http.ResponseWriter, r * http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}

		createConnection(ws)
	}
}