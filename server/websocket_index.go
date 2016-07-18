package server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/poslegm/blockchain-chat/db"
	"github.com/poslegm/blockchain-chat/message"
	"github.com/poslegm/blockchain-chat/network"
	"net/http"
)

// выбирает ответ на сообщение в зависимости от типа и кладёт его в очередь сообщений
func switchTypes(msg WebSocketMessage) {
	fmt.Println("Websockets.swithTypes: ", msg)

	switch msg.Type {
	case "GetMessages":
		sendAllMessages()
	case "SendMessage":
		sendMessageToNetwork(msg)
	case "GetMyKey":
		sendPublicKey()
	case "GetContacts":
		sendContacts()
	}
}

func sendAllMessages() {
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
			textMsg.Receiver, textMsg.Sender, textMsg.Text, false,
		})
	}

	WebSocketQueue <- WebSocketMessage{Type: "AllMessages", Messages: chatMessages}
}

func sendMessageToNetwork(msg WebSocketMessage) {
	if len(msg.Messages) != 1 {
		fmt.Printf("WebSocket.switchTypes: incorrect message - %#v\n", msg)
		return
	}
	chatMsg := msg.Messages[0]

	if chatMsg.NewPublicKey {
		err := chatMsg.addNewPublicKeyToDb()
		if err != nil {
			fmt.Println("Websocket.swithTypes: can't add new public key ", err.Error())
		}
	}

	kp, err := db.GetKeyByAddress(chatMsg.Receiver)

	if chatMsg.NewPublicKey {
		kp = &message.KeyPair{[]byte(chatMsg.Receiver), []byte{}, []byte{}}
	}

	if err != nil {
		fmt.Println("WebSockets.swithTypes: can't get kp from db ", err.Error())
		return
	} else if kp == nil {
		fmt.Println("WebSockets.swithTypes: there is no kp in db")
		return
	}

	networkMsg, err := network.CreateTextNetworkMessage(
		kp.GetBase58Address(),
		chatMsg.Sender,
		chatMsg.Text,
		kp.PublicKey,
	)

	if err != nil {
		fmt.Println("Websockets.switchTypes: can't send message ", err.Error())
	} else {
		go network.CurrentNetworkUser.SendMessage(networkMsg)
	}

	if chatMsg.NewPublicKey {
		WebSocketQueue <- WebSocketMessage{
			Type:     "NewKeyHash",
			Key:      kp.GetBase58Address(),
			Messages: []ChatMessage{chatMsg},
		}
	}
}

func sendPublicKey() {
	publicKey, err := db.GetPublicKey()
	if err != nil {
		fmt.Println("Websockets.switchTypes: can't send public key ", err.Error())
	} else {
		WebSocketQueue <- WebSocketMessage{Type: "Key", Key: string(publicKey)}
	}
}

func sendContacts() {
	contacts, err := db.GetAllContacts()
	if err != nil {
		fmt.Println("Websockets.sendContacts: can't get contacts from db ", err.Error())
		return
	}

	contactsMessage := make([]string, len(contacts))
	for i, contact := range contacts {
		contactsMessage[i] = contact.GetBase58Address()
	}

	WebSocketQueue <- WebSocketMessage{Type: "AllContacts", Contacts: contactsMessage}
}

func createWSHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
		if err != nil {
			http.Error(w, "Bad request", 400)
			return
		}

		createConnection(ws, switchTypes)
	}
}
