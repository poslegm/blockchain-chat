package db

import (
	"testing"
	"github.com/poslegm/blockchain-chat/network"
	"fmt"
	"time"
	"github.com/poslegm/blockchain-chat/message"
)

func TestDB(t *testing.T) {
	err := InitDB()
	if err != nil {
		t.Fatal(err)
	}
	err = CloseDB()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMessages(t *testing.T) {
	tInitDB()
	defer tCloseDB()

	inMsgs := []network.NetworkMessage{
		//network.CreateTextNetworkMessage("123", "345", "678"),
		//network.CreateTextNetworkMessage("asd", "fgh", "jkl"),
		//network.CreateTextNetworkMessage("zxc", "vbn", "nmm"),
	}

	err := AddMessages(inMsgs)
	if err != nil {
		t.Fatal(err)
	}

	outMsgs, err := GetAllMessages()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range outMsgs {
		fmt.Println(v)
	}
}

func TestAddresses(t *testing.T) {
	tInitDB()
	defer tCloseDB()

	inAddresses := []network.NetAddress{
		network.CreateNetAddress(time.Now(), "123.123.123", "10"),
		network.CreateNetAddress(time.Now(), "213.213.213", "11"),
		network.CreateNetAddress(time.Now(), "143.143.143", "12"),
	}

	err := AddKnownAddresses(inAddresses)
	if err != nil {
		t.Fatal(err)
	}

	outAddresses, err := GetKnownAddresses()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range outAddresses {
		fmt.Println(v)
	}
}

func TestKeys(t *testing.T) {
	tInitDB()
	defer tCloseDB()

	inKeys := []*message.KeyPair{
		&message.KeyPair{[]byte("pubkey1"), []byte("privkey1"), []byte("passphrase")},
		&message.KeyPair{[]byte("asdfkas"), []byte("1234234"), []byte("sdnfsj")},
	}
	for _, v := range inKeys {
		fmt.Println(v)
	}
	err := AddKeys(inKeys)
	if err != nil {
		t.Fatal(err)
	}

	outKeys, err := GetAllKeys()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range outKeys {
		fmt.Println(v)
	}

	addr := inKeys[0].GetBase58Address()
	fmt.Println(addr)

	outKey, err := GetKeyByAddress(addr)
	if err != nil {
		t.Fatal(err)
	}
	if outKey == nil {
		t.Fatal("key wasn't found")
	}
	fmt.Println(outKey)
}
