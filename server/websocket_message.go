package server

import (
	"github.com/poslegm/blockchain-chat/db"
	"github.com/poslegm/blockchain-chat/message"
)

type WebSocketMessage struct {
	Type     string
	Messages []ChatMessage
	Key      string
}

type ChatMessage struct {
	Receiver     string
	Sender       string
	Text         string
	NewPublicKey bool
}

func (msg ChatMessage) addNewPublicKeyToDb() error {
	return db.AddContacts([]*message.KeyPair{{
		PublicKey: []byte(msg.Receiver),
		PrivateKey: []byte{},
		Passphrase: []byte{},
	}})
}