package network

import (
	"crypto/md5"
	"encoding/json"

	"errors"
	"fmt"
	"github.com/poslegm/blockchain-chat/message"
	"github.com/poslegm/blockchain-chat/shahash"
)

const MESSAGE = "MESSAGE"
const REQUEST = "REQUEST"

type NetworkMessage struct {
	MessageType string
	// заполняется, если messageType == REQUEST, содержит адрес, к которому пытается подключиться
	IP string

	Data []byte
}

var Hash = md5.Sum

func CreateTextNetworkMessage(receiver, sender, text string,
	time int64, publicKey []byte, parent *message.TextMessage) (NetworkMessage, error) {
	var (
		parentHash shahash.ShaHash = shahash.ShaHash{}
		err        error           = nil
	)
	if parent != nil {
		parentHash, err = message.GenerateParentHash(*parent)
		if err != nil {
			fmt.Println("network_message.CreateTextNetworkMessage: can't generate parent hash ", err.Error())
			return NetworkMessage{}, err
		}
	}

	textMessage := message.TextMessage{
		Receiver:   receiver,
		Sender:     sender,
		Text:       text,
		Time:       time,
		ParentHash: parentHash,
	}

	err = textMessage.Mine()
	if err != nil {
		fmt.Println("network_message.CreateTextNetworkMessage: can't mine ", err.Error())
		return NetworkMessage{}, err
	}
	fmt.Println("PRE-SEND: ", textMessage)
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
			txtMsg, err := encrypted.Decode(kp)
			if err != nil {
				return message.TextMessage{}, err
			} else {
				return verify(txtMsg)
			}
		}
	}
	return message.TextMessage{}, errors.New("unsuitable-pair")
}

func verify(textMsg message.TextMessage) (message.TextMessage, error) {
	verified, err := textMsg.Verify()
	if verified {
		fmt.Println("VERIFIED")
		return textMsg, err
	} else {
		fmt.Println("NOT VERIFIED")
		return message.TextMessage{}, err
	}
}
