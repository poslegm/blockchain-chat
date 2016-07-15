package server

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/poslegm/blockchain-chat/db"
	"github.com/poslegm/blockchain-chat/message"
	"fmt"
)

func handleAddition(msg WebSocketMessage) {
	switch msg.Type {
	case "Contact":
		kp := &message.KeyPair{[]byte(msg.Key), []byte{}, []byte{}}
		err := db.AddContacts([]*message.KeyPair{kp})
		if err != nil {
			fmt.Println("WebsocketAddition.handleAddition: can't add contact ", err.Error())
			WebSocketQueue <- WebSocketMessage{Type: "BadContact"}
		} else {
			WebSocketQueue <- WebSocketMessage{Type: "OkContact", Key: kp.GetBase58Address()}
		}
	case "KeyPair":
		kp := &message.KeyPair{
			PublicKey:  []byte(msg.PublicKey),
			PrivateKey: []byte(msg.PrivateKey),
			Passphrase: []byte(msg.Passphrase),
		}
		err := db.AddKeys([]*message.KeyPair{kp})
		if err != nil {
			fmt.Println("WebsocketAddition.handleAddition: can't add key pair ", err.Error())
			WebSocketQueue <- WebSocketMessage{Type: "BadKeyPair"}
		} else {
			WebSocketQueue <- WebSocketMessage{Type: "OkKeyPair"}
		}
	}
}

func createAdditionWSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}

		createConnection(ws, handleAddition)
	}
}
