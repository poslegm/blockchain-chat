package network

import (
	"encoding/json"

	"errors"
	"fmt"
	"github.com/poslegm/blockchain-chat/message"
)

const MESSAGE = "MESSAGE"
const REQUEST = "REQUEST"

type NetworkMessage struct {
	MessageType string
	// заполняется, если messageType == REQUEST, содержит адрес, к которому пытается подключиться
	IP string

	Data []byte
}

func CreateTextNetworkMessage(receiver, sender, text string, time int64, publicKey []byte) (NetworkMessage, error) {
	textMessage := message.TextMessage{
		Receiver: receiver,
		Sender:   sender,
		Text:     text,
		Time:     time,
	}

	encrypted, err := textMessage.Encode(&message.KeyPair{publicKey, []byte{}, []byte{}})
	if err != nil {
		fmt.Println("network_message.CreateTextNetworkMessage: ", err.Error())
		return NetworkMessage{}, err
	}

	encryptedBytes, err := json.Marshal(encrypted)

	return NetworkMessage{MessageType: MESSAGE, Data: encryptedBytes}, err
}

func (msg NetworkMessage) AsTextMessage() (message.TextMessage, error) {
	if msg.MessageType != MESSAGE {
		return message.TextMessage{}, errors.New(
			"network_message.asTextMessage: not encrypted messages " + msg.MessageType)
	}
	encrypted := message.EncryptedMessage{}
	err := json.Unmarshal(msg.Data, &encrypted)
	if err != nil {
		return message.TextMessage{}, err
	}

	for _, kp := range CurrentNetworkUser.KeyPairs {
		if encrypted.ReceiverAddress == kp.GetBase58Address() {
			fmt.Println("DECODING...")
			return encrypted.Decode(kp)
		}
	}
	return message.TextMessage{}, errors.New("unsuitable-pair")
}
