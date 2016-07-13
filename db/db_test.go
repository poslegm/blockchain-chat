package db

import (
	"testing"
	"github.com/poslegm/blockchain-chat/network"
	"fmt"
)

func TestDB(t *testing.T) {
	InitDB()
	msgs := []network.NetworkMessage{
		network.NetworkMessage{},
		network.NetworkMessage{},
		network.NetworkMessage{},
	}
	addrs := []network.NetAddress{
		network.NetAddress{},
		network.NetAddress{},
		network.NetAddress{},
	}
	AddKnownAddresses(addrs)
	AddMessages(msgs)
	gm, err := GetAllMessages()
	if(err != nil) {
		t.Errorf("%s", err)
	}
	for _, v := range gm {
		fmt.Println(v)
	}
	ga, err := GetKnownAddresses()
	if(err != nil) {
		t.Errorf("%s", err)
	}
	for _, v := range ga {
		fmt.Println(v)
	}
}
