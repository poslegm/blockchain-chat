package server

import (
	"errors"
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
		fmt.Println("Websockets.switchTypes: can't get network messages ", err.Error())
		return
	}
	textMessages, err := db.GetAllTextMessages()
	if err != nil {
		fmt.Println("Websockets.switchTypes: can't get text messages ", err.Error())
		return
	}

	// отправляются свои сообщения и те, которые можно расшифровать
	chatMessages := make([]ChatMessage, len(textMessages))
	for i, textMsg := range textMessages {
		chatMessages[i] = ChatMessage{
			textMsg.Receiver, textMsg.Sender, textMsg.Text, textMsg.Time, false,
		}
	}

	for _, networkMsg := range networkMessages {
		textMsg, err := networkMsg.AsTextMessage()
		if err != nil {
			if err.Error() != "unsuitable-pair" {
				fmt.Println("Websockets.switchTypes: ", err.Error())
			}
			continue
		}

		chatMessages = append(chatMessages, ChatMessage{
			textMsg.Receiver, textMsg.Sender, textMsg.Text, textMsg.Time, false,
		})
	}

	WebSocketQueue <- WebSocketMessage{Type: "AllMessages", Messages: chatMessages}
}

func sendMessageToNetwork(msg WebSocketMessage) {
	if len(msg.Messages) != 1 {
		fmt.Printf("websockets_index.sendMessageToNetwork: incorrect message - %#v\n", msg)
		return
	}
	chatMsg := msg.Messages[0]

	if chatMsg.NewPublicKey {
		err := chatMsg.addNewPublicKeyToDb()
		if err != nil {
			fmt.Println("websockets_index.sendMessageToNetwork: can't add new public key ", err.Error())
		}
	}

	kp, err := db.GetContactByAddress(chatMsg.Receiver)

	if chatMsg.NewPublicKey {
		kp = &message.KeyPair{[]byte(chatMsg.Receiver), []byte{}, []byte{}}
	}

	if err != nil {
		fmt.Println("websockets_index.sendMessageToNetwork: can't get kp from db ", err.Error())
		return
	} else if kp == nil {
		fmt.Println("websockets_index.sendMessageToNetwork: there is no kp in db")
		return
	}

	parent, err := getLastMessageInDialog(chatMsg.Receiver, chatMsg.Sender)
	if err != nil {
		fmt.Println("websockets_index.sendMessageToNetwork: can't get parent ", err.Error())
	}

	networkMsg, err := network.CreateTextNetworkMessage(
		kp.GetBase58Address(),
		chatMsg.Sender,
		chatMsg.Text,
		chatMsg.Time,
		kp.PublicKey,
		parent,
	)

	if err != nil {
		fmt.Println("websockets_index.sendMessageToNetwork: can't send message ", err.Error())
	} else {
		network.CurrentNetworkUser.OutgoingMessages <- networkMsg
		go network.CurrentNetworkUser.SendMessage(networkMsg)
	}

	if chatMsg.NewPublicKey {
		WebSocketQueue <- WebSocketMessage{
			Type:     "NewKeyHash",
			Key:      kp.GetBase58Address(),
			Messages: []ChatMessage{chatMsg},
		}
	}

	saveMyMessage(chatMsg, kp)
}

func getLastMessageInDialog(talker1 string, talker2 string) (*message.TextMessage, error) {
	last1, err := db.GetLastTextMessageFromSender(talker1)
	if err != nil {
		fmt.Println("websocket_index.getLastMessageInDialog: can't get last msg ", err.Error())
		return nil, err
	}
	last2, err := db.GetLastTextMessageFromSender(talker2)
	if err != nil {
		fmt.Println("websocket_index.getLastMessageInDialog: can't get last msg ", err.Error())
		return nil, err
	}

	if last1.Time == 0 && last2.Time == 0 {
		return nil, errors.New("No parent")
	}

	if last1.Time > last2.Time {
		return &last1, nil
	} else {
		return &last2, nil
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

func saveMyMessage(chatMsg ChatMessage, kp *message.KeyPair) {
	err := db.AddTextMessages([]message.TextMessage{{
		Receiver: kp.GetBase58Address(),
		Sender:   chatMsg.Sender,
		Text:     chatMsg.Text,
		Time:     chatMsg.Time,
	}})

	if err != nil {
		fmt.Println("websockets.SaveMyMessage: can't save message ", err.Error())
	}
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
