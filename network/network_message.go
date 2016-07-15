package network

import (
	"encoding/json"

	"github.com/poslegm/blockchain-chat/message"
	"errors"
	"fmt"
)

// TODO сообщение, идущее по сети, должно быть сериализованно в массив байтов;
// TODO сделать интерфейс для преобразования сообщения от клиента в вид для сети;
// TODO так же по сети могут передаваться блоки, надо держать это в голове;

const MESSAGE = "MESSAGE"
const REQUEST = "REQUEST"

type NetworkMessage struct {
	MessageType string
	// заполняется, если messageType == REQUEST, содержит адрес, к которому пытается подключиться
	IP          string

	Data        []byte
}

func CreateTextNetworkMessage(receiver, sender, text string, publicKey []byte) (NetworkMessage, error) {
	textMessage := message.TextMessage{
		Receiver:receiver,
		Sender:sender,
		Text:text,
	}

	encrypted, err := textMessage.Encode(&message.KeyPair{publicKey, []byte{}, []byte{}})
	if err != nil {
		fmt.Println("network_message.CreateTextNetworkMessage: ", err.Error())
		return NetworkMessage{}, err
	}

	encryptedBytes, err := json.Marshal(encrypted)

	return NetworkMessage{MessageType:MESSAGE, Data:encryptedBytes}, err
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
			return encrypted.Decode(kp)
		}
	}
	return message.TextMessage{}, errors.New("unsuitable-pair")
}