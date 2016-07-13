package message

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/maxwellhealth/go-gpg"
	"bytes"
	"fmt"
	"io/ioutil"
)

//gpg key pair
type KeyPair struct {
	//gpg pub key
	PublicKey []byte

	//gpg private key
	PrivateKey []byte
}

//encode data using kp's public key
func (kp KeyPair) Encode(data []byte) ([]byte, error) {
	//check publickey existence
	if kp.PublicKey == nil {
		return nil, fmt.Errorf("no public key provided")
	}

	//create buffers
	inputBuffer := bytes.NewBuffer(data)
	var outputBuffer bytes.Buffer

	//encode
	err := gpg.Encode(kp.PublicKey, inputBuffer, &outputBuffer)
	if err != nil {
		return nil, fmt.Errorf("error encoding data: %s", err)
	}

	return ioutil.ReadAll(&outputBuffer)
}

//decode data using kp's private key
func (kp KeyPair) Decode(data []byte) ([]byte, error) {
	//check privatekey existence
	if kp.PrivateKey == nil {
		return nil, fmt.Errorf("no private key provided")
	}

	//create buffers
 	inputBuffer := bytes.NewBuffer(data)
	var outputBuffer bytes.Buffer

	//decode
	err := gpg.Decode(kp.PrivateKey, []byte{}, inputBuffer, &outputBuffer)
	if err != nil {
		return nil, fmt.Errorf("error decoding data: %s", err)
	}

	return ioutil.ReadAll(&outputBuffer)
}

//get address to send to
func (kp KeyPair) GetBase58Address() string {
	return base58.Encode(kp.PublicKey)
}

//get public key from address
func KeyPairFromBase58Address(address string) KeyPair {
	return KeyPair{PublicKey:base58.Decode(address)}
}