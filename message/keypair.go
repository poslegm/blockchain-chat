package message

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/maxwellhealth/go-gpg"
	"bytes"
	"fmt"
	"io/ioutil"
	"golang.org/x/crypto/openpgp"
	"crypto/md5"
)

//gpg key pair
type KeyPair struct {
	//gpg pub key
	PublicKey  []byte

	//gpg private key
	PrivateKey []byte

	//gpg private key passphrase
	Passphrase []byte
}

//encode data using kp's public key
func (kp *KeyPair) Encode(data []byte) ([]byte, error) {
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
func (kp *KeyPair) Decode(data []byte) ([]byte, error) {
	//check privatekey existence
	if kp.PrivateKey == nil {
		return nil, fmt.Errorf("no private key provided")
	}

	//create buffers
 	inputBuffer := bytes.NewBuffer(data)
	var outputBuffer bytes.Buffer

	//decode
	err := gpg.Decode(kp.PrivateKey, kp.Passphrase, inputBuffer, &outputBuffer)
	if err != nil {
		return nil, fmt.Errorf("error decoding data: %s", err)
	}

	return ioutil.ReadAll(&outputBuffer)
}

//get address to send to
func (kp *KeyPair) GetBase58Address() string {
	sum := md5.Sum(kp.PublicKey)
	return base58.Encode(sum[:])
}

//string representation
func (kp * KeyPair) String() string {
	return "pub:" + string(kp.PublicKey) + "\npriv:" + string(kp.PrivateKey) + "\npassphrase:" + string(kp.Passphrase)
}

func KeyPairFromFile(publicKeyFile, privateKeyFile, passphrase string) (*KeyPair, error) {
	pub, err := ioutil.ReadFile(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("keypair from file cannot open %s: %s", publicKeyFile, err)
	}
	priv, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return nil, fmt.Errorf("keypair from file cannot open %s: %s", privateKeyFile, err)
	}
	_, err = openpgp.ReadArmoredKeyRing(bytes.NewBuffer(pub))
	if err != nil {
		return nil, fmt.Errorf("public key error: %s", err)
	}
	_, err = openpgp.ReadArmoredKeyRing(bytes.NewBuffer(priv))
	if err != nil {
		return nil, fmt.Errorf("private key error: %s", err)
	}
	return &KeyPair{PrivateKey:priv, PublicKey:pub, Passphrase:[]byte(passphrase)}, nil
}