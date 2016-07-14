package server

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
	"time"
	"github.com/poslegm/blockchain-chat/network"
	"github.com/poslegm/blockchain-chat/db"
)

var WebSocketQueue = make(chan WebSocketMessage)

type WebSocketMessage struct {
	Type     string
	Messages []ChatMessage
	Key string
}

type ChatMessage struct {
	Receiver string
	Sender string
	Text string
}

func receive(ws *websocket.Conn) {
	// чтение не должно прекращаться
	ws.SetReadDeadline(time.Time{})

	for {
		fmt.Println("KEK")
		msg := WebSocketMessage{}

		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("WebSockets.receive: " + err.Error())
			break
		}

		msg.switchTypes()
	}
	fmt.Println("Выход из цикла")
	ws.Close()
}

// выбирает ответ на сообщение в зависимости от типа и кладёт его в очередь сообщений
func (msg WebSocketMessage) switchTypes() {
	fmt.Println("Websockets.swithTypes: ", msg)

	switch msg.Type {
	case "GetMessages":
		// TODO выводить только сообщения, которые можно расшифровать
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

		WebSocketQueue <- WebSocketMessage{Type:"AllMessages", Messages:chatMessages}
	case "SendMessage":
		if len(msg.Messages) != 1 {
			fmt.Printf("WebSocket.switchTypes: incorrect message - %#v\n", msg)
			return
		}
		chatMsg := msg.Messages[0]
		go network.CurrentNetworkUser.SendMessage(network.CreateTextMessage(
			chatMsg.Receiver,
			chatMsg.Sender,
			chatMsg.Text,
		))
	case "GetMyKey":
		WebSocketQueue <- WebSocketMessage{Type:"Key", Key:db.GetPublicKey()}
	}
}

func handleMessagesQueue(ws *websocket.Conn) {
	for {
		select {
		case msg := <- WebSocketQueue:
			ws.SetWriteDeadline(time.Now().Add(30 * time.Second))
			err := ws.WriteJSON(&msg)
			if err != nil {
				fmt.Println("WebSocket.handleMessagesQueue: " + err.Error())
				// если сообщение не удалось отправить, оно добавляется обратно в конец очереди
				WebSocketQueue <- msg
			} else {
				fmt.Println("WebSocket.handleMessagesQueue: sended ", msg)
			}
		case msg := <- network.CurrentNetworkUser.IncomingMessages:
			WebSocketQueue <- WebSocketMessage{
				Type:"NewMessage",
				Messages:[]ChatMessage{{
					msg.Receiver,
					msg.Sender,
					msg.Text,
				}},
			}
		}
	}
}

func createConnection(ws *websocket.Conn) {
	go receive(ws)
	go handleMessagesQueue(ws)
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