package message

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/poslegm/blockchain-chat/shahash"
)

type EncryptedMessage struct {
	//base58 encoded hash of public key
	ReceiverAddress string

	//encrypted message contents
	DataLength int
	Data       []byte
}

type TextMessage struct {
	Receiver string
	Sender   string
	Text     string
	Time     int64

	ParentHash  shahash.ShaHash
	MessageHash shahash.ShaHash

	Height int64
	Nonce  int64
}

func itob(v int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v))
	return buf
}

func (msg TextMessage) toByteArray() []byte {
	buf := []byte(msg.Receiver + msg.Sender + msg.Text)
	buf = append(buf, itob(msg.Time)...)
	buf = append(buf, msg.ParentHash[:]...)
	buf = append(buf, itob(msg.Height)...)
	buf = append(buf, itob(msg.Nonce)...)
	return buf
}

func (msg *TextMessage) Mine() (err error) {
	buf := msg.toByteArray()
	msg.MessageHash, err = shahash.ShaHashFromData(buf)
	if err != nil {
		return fmt.Errorf("mining error: %s", err)
	}
	fmt.Println("Mining...")
	for !msg.MessageHash.Check() {
		msg.Nonce++
		buf = append(buf[:len(buf)-8], itob(msg.Nonce)...)
		msg.MessageHash, err = shahash.ShaHashFromData(buf)
		if err != nil {
			return fmt.Errorf("mining error: %s", err)
		}
	}
	return nil
}

func GenerateParentHash(parent TextMessage) (shahash.ShaHash, error) {
	return shahash.ShaHashFromData(parent.toByteArray())
}

func (msg TextMessage) Verify() (bool, error) {
	buf := msg.toByteArray()
	hash, err := shahash.ShaHashFromData(buf)
	if err != nil {
		return false, fmt.Errorf("verify hash func error: %s", err)
	}

	if !hash.Equal(msg.MessageHash) {
		return false, nil
	}
	return true, nil
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

	return EncryptedMessage{
		ReceiverAddress: kp.GetBase58Address(),
		DataLength:      len(encoded),
		Data:            encoded,
	}, nil
}
