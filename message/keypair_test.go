package message

import (
	"testing"
	"fmt"
	"io/ioutil"
	"os"
	_ "golang.org/x/crypto/ripemd160"
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
		t.Fatalf("cannot create keypair from file: %s", err)
	}
	encoded, err := kp.Encode([]byte(message))
	ioutil.WriteFile("byte.msg", encoded, 0660)
	if err != nil {
		t.Fatalf("kp encode error: %s", err)
	}
	//fmt.Println(string(encoded))
	decoded, err := kp.Decode(encoded)
	if err != nil {
		t.Fatalf("kp decode error: %s", err)
	}
	fmt.Println(string(decoded))
}

func TestGenerateKeyPair(t *testing.T) {
	kp, err := GenerateKeyPair("asd@mail", "lol", "kek", "")
	if err != nil {
		t.Fatalf("generate key pair: %s", err)
	}

	err = kp.SaveToFile("new-key")
	if err != nil {
		t.Fatalf("save to file: %s", err)
	}

	nkp, err := KeyPairFromFile("new-key.pub", "new-key.priv", "new-key.pass")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nkp)

	encoded, err := nkp.Encode([]byte("hello world!"))
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := nkp.Decode(encoded)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(decoded))
	os.Remove("new-key.pub")
	os.Remove("new-key.priv")
	os.Remove("new-key.pass")
}