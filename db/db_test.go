package db

import (
	"testing"
	"github.com/poslegm/blockchain-chat/network"
	"fmt"
	"time"
	"net"
)

func TestDB(t *testing.T) {
	InitDB()
	msgs := []network.NetworkMessage{
		network.NetworkMessage{"1", "2", "3"},
		network.NetworkMessage{"4", "5", "6"},
		network.NetworkMessage{"7", "8", "9"},
	}
	addrs := []network.NetAddress{
		network.NetAddress{time.Time{}, net.IPv4(1, 2, 3, 4), 10},
		network.NetAddress{time.Time{}, net.IPv4(1, 2, 3, 4), 15},
		network.NetAddress{time.Time{}, net.IPv4(1, 2, 3, 4), 15},
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
