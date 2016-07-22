package db

import (
	"testing"
	"github.com/poslegm/blockchain-chat/network"
	"fmt"
	"time"
	"github.com/poslegm/blockchain-chat/message"
)

func TestDB(t *testing.T) {
	err := tInitDB()
	if err != nil {
		t.Fatal(err)
	}
	err = tCloseDB()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMessages(t *testing.T) {
	tInitDB()
	defer tCloseDB()

	inMsgs := []network.NetworkMessage{
		{"123", "asd", nil},
		{"asd", "get", nil},
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

	has, err := HasMessage(inMsgs[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(has)
	err = AddMessages(inMsgs[0:0])
	if err != nil {
		t.Fatal(err)
	}

	outMsgs, err = GetAllMessages()
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

	inAddresses[0].Ip = "123.234.123"
	err = AddKnownAddresses(inAddresses)
	if err != nil {
		t.Fatal(err)
	}

	outAddresses, err = GetKnownAddresses()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println()
	for _, v := range outAddresses {
		fmt.Println(v)
	}
}

func TestKeys(t *testing.T) {
	tInitDB()
	defer tCloseDB()

	inKeys := []*message.KeyPair{
		{[]byte("pubkey1"), []byte("privkey1"), []byte("passphrase")},
		{[]byte("asdfkas"), []byte("1234234"), []byte("sdnfsj")},
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

func TestContacts(t *testing.T) {
	tInitDB()
	defer tCloseDB()

	inContacts := []*message.KeyPair{
		{[]byte("pubkey1"), nil, nil},
		{[]byte("pubkey2"), nil, nil},
	}

	for _, v := range inContacts {
		fmt.Println(v)
	}

	err := AddContacts(inContacts)
	if err != nil {
		t.Fatal(err)
	}

	outContacts, err := GetAllContacts()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range outContacts {
		fmt.Println(v)
	}

	addr := inContacts[0].GetBase58Address()
	fmt.Println(addr)

	outContact, err := GetContactByAddress(addr)
	if err != nil {
		t.Fatal(err)
	}
	if outContact == nil {
		t.Fatal("contact wasn't found")
	}
	fmt.Println(outContact)
}

func makeTextMessages() []message.TextMessage {
	msgs := []message.TextMessage{{}, {}, {}}
	msgs[0].Sender = "send1"
	msgs[1].Sender = "send2"
	msgs[2].Sender = "send3"
	return msgs
}

func TestTextMessages(t *testing.T) {
	tInitDB()
	defer tCloseDB()

	inMessages := makeTextMessages()
	for i := 0; i < len(inMessages); i++ {
		inMessages[i].Mine()
	}
	err := AddTextMessages(inMessages)
	if err != nil {
		t.Fatal(err)
	}

	outMessages, err := GetAllTextMessages()
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range outMessages {
		ver, err := v.Verify()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(v, " ", ver)
	}
	for i := 0; i < len(inMessages); i++ {
		msg, err := GetTextMessagesBySender(inMessages[len(inMessages) - 1 - i].Sender)
		if err != nil {
			t.Fatal(err)
		}
		ver, err := msg[0].Verify()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(msg, " ", ver, "\n")
	}
	msg := message.TextMessage{}
	msg.Sender = inMessages[0].Sender
	msg.Text = "last"
	err = AddTextMessages([]message.TextMessage{msg})
	if err != nil {
		t.Fatal(err)
	}
	msg, err = GetLastTextMessageFromSender(inMessages[0].Sender)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(msg)
}
