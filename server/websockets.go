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
		msg := WebSocketMessage{}

		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("WebSockets.receive: " + err.Error())
			break
		}

		msg.switchTypes()
	}
	ws.Close()
}

// выбирает ответ на сообщение в зависимости от типа и кладёт его в очередь сообщений
func (msg WebSocketMessage) switchTypes() {
	fmt.Println("Websockets.swithTypes: ", msg)

	switch msg.Type {
	case "GetMessages":
		networkMessages, err := db.GetAllMessages()
		if err != nil {
			fmt.Println("Websockets.switchTypes: ", err.Error())
			return
		}
		fmt.Println(networkMessages)
		chatMessages := make([]ChatMessage, 0)
		for _, networkMsg := range networkMessages {
			textMsg, err := networkMsg.AsTextMessage()
			if err != nil {
				if err.Error() != "unsuitable-pair" {
					fmt.Println("Websockets.switchTypes: ", err.Error())
				}
				continue
			}

			chatMessages = append(chatMessages, ChatMessage{
				textMsg.Receiver, textMsg.Sender, textMsg.Text,
			})
		}

		WebSocketQueue <- WebSocketMessage{Type:"AllMessages", Messages:chatMessages}
	case "SendMessage":
		if len(msg.Messages) != 1 {
			fmt.Printf("WebSocket.switchTypes: incorrect message - %#v\n", msg)
			return
		}
		chatMsg := msg.Messages[0]

		kp, err := db.GetKeyByAddress(chatMsg.Receiver)
		if err != nil || kp == nil {
			fmt.Println("WebSockets.swithTypes: can't get kp from db ", err.Error())
			return
		}
		networkMsg, err := network.CreateTextNetworkMessage(
			chatMsg.Receiver,
			chatMsg.Sender,
			chatMsg.Text,
			kp.PublicKey,
		)

		if err != nil {
			fmt.Println("Websockets.switchTypes: can't send message ", err.Error())
		} else {
			go network.CurrentNetworkUser.SendMessage(networkMsg)
		}
	case "GetMyKey":
		publicKey, err := db.GetPublicKey()
		if err != nil {
			fmt.Println("Websockets.switchTypes: can't send public key ", err.Error())
		} else {
			WebSocketQueue <- WebSocketMessage{Type:"Key", Key:string(publicKey)}
		}
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
			textMsg, err := msg.AsTextMessage()
			if err != nil {
				if err.Error() != "unsuitable-pair" {
					fmt.Println("Websockts.switchTypes: ", err.Error())
				}
				continue
			} else {
				WebSocketQueue <- WebSocketMessage{
					Type:"NewMessage",
					Messages:[]ChatMessage{{
						textMsg.Receiver,
						textMsg.Sender,
						textMsg.Text,
					}},
				}
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