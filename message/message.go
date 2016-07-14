package message

import (
	"encoding/json"
)

type EncryptedMessage struct {
	//base58 encoded hash of public key
	ReceiverAddress string

	//encrypted message contents
	DataLength int
	Data []byte
}

type TextMessage struct {
	Receiver string
	Sender string
	Text string
}

func (msg EncryptedMessage) Decode(kp *KeyPair) (TextMessage, error) {
	encrypted, err := kp.Decode(msg.Data)
	if err != nil {
		return TextMessage{}, err
	}

	textMessage := TextMessage{}
	err = json.Unmarshal(encrypted, &textMessage)
	return textMessage, err
}

func (msg TextMessage) Encode(kp *KeyPair) (EncryptedMessage, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return EncryptedMessage{}, err
	}

	encoded, err := kp.Encode(data)
	if err != nil {
		return EncryptedMessage{}, err
	}

	return EncryptedMessage {
		ReceiverAddress:kp.GetBase58Address(),
		DataLength:len(encoded),
		Data:encoded,
	}, nil
}