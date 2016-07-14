package message

type EncryptedMessage struct {
	//base58 encoded hash of public key
	ReceiverAddress string

	//encrypted message contents
	DataLength int
	Data []byte
}


