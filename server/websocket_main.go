package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/poslegm/blockchain-chat/network"
	"time"
)

var WebSocketQueue = make(chan WebSocketMessage)

func receive(ws *websocket.Conn, handle func(m WebSocketMessage)) {
	// чтение не должно прекращаться
	ws.SetReadDeadline(time.Time{})

	for {
		msg := WebSocketMessage{}

		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("WebSockets.receive: " + err.Error())
			break
		}

		handle(msg)
	}
	ws.Close()
}

func handleMessagesQueue(ws *websocket.Conn) {
	for {
		select {
		case msg := <-WebSocketQueue:
			writeMessageToWS(msg, ws)
		case msg := <-network.CurrentNetworkUser.IncomingMessages:
			writeNetworkMessageToQueue(msg)
		}
	}
}

func writeMessageToWS(msg WebSocketMessage, ws *websocket.Conn) {
	ws.SetWriteDeadline(time.Now().Add(30 * time.Second))
	err := ws.WriteJSON(&msg)
	if err != nil {
		fmt.Println("WebSocket.handleMessagesQueue: " + err.Error())
		// если сообщение не удалось отправить, оно добавляется обратно в конец очереди
		WebSocketQueue <- msg
	} else {
		fmt.Println("WebSocket.handleMessagesQueue: sended ", msg)
	}
}

func writeNetworkMessageToQueue(msg network.NetworkMessage) {
	textMsg, err := msg.AsTextMessage()
	if err != nil {
		if err.Error() != "unsuitable-pair" {
			fmt.Println("Websockts.switchTypes: ", err.Error())
		}
		return
	} else {
		WebSocketQueue <- WebSocketMessage{
			Type: "NewMessage",
			Messages: []ChatMessage{{
				Receiver:     textMsg.Receiver,
				Sender:       textMsg.Sender,
				Text:         textMsg.Text,
				Time:         textMsg.Time,
				NewPublicKey: false,
			}},
		}
	}
}

func createConnection(ws *websocket.Conn, handle func(m WebSocketMessage)) {
	go receive(ws, handle)
	go handleMessagesQueue(ws)
}
