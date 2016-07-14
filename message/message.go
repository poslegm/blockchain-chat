package message

type Address string

type EncryptedMessage struct {
	//base58 encode pubkey
	ReceiverAddress Address

	//encrypted message contents
	DataLength int
	Data []byte
}



