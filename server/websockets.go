package server

import (
	"net/http"
	"github.com/gorilla/websocket"
	"fmt"
)

// сообщения для веб-сокета:
// запрос на получение всех сообщений (которые можно расшифровать)
// запрос на отправку сообщения (набранный текст + ключ получателя)
// запрос на состояние майнинга (решить, что какое там может быть состояние; время?)

type WebSocketMessage struct {
	Command string
	Value string
}

func receive(ws *websocket.Conn) {
	// чтение не должно прекращаться
	ws.SetReadDeadline(0)

	defer ws.Close()
	for {
		msg := WebSocketMessage{}

		err := ws.ReadJSON(&msg)
		if err != nil {
			fmt.Println("WebSockets.receive json: " + err.Error())
			return
		}

		fmt.Printf("WebSocket.receive: message %#v\n", msg)
	}
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