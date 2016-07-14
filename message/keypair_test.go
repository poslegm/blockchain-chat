package message

import (
	"testing"
	"fmt"
	"io/ioutil"
)

const (
	pubKey = "samplekey_pub.asc"
	privKey = "samplekey_priv.asc"
	passphrase = "sample-key"
	message = "sample key pair test yeah!"
)

func TestKeyPair(t *testing.T) {
	kp, err := KeyPairFromFile(pubKey, privKey, passphrase)
	if err != nil {
		t.Errorf("cannot create keypair from file: %s", err)
		t.FailNow()
	}
	encoded, err := kp.Encode([]byte(message))
	ioutil.WriteFile("byte.msg", encoded, 0660)
	if err != nil {
		t.Errorf("kp encode error: %s", err)
		t.FailNow()
	}
	//fmt.Println(string(encoded))
	decoded, err := kp.Decode(encoded)
	if err != nil {
		t.Errorf("kp decode error: %s", err)
		t.FailNow()
	}
	fmt.Println(string(decoded))
}

func TestBase58(t *testing.T) {

}