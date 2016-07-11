package network

import "fmt"

type NetworkMessage struct {
	Receiver string
	Sender string
	Text string
}

func SendMessage(msg NetworkMessage) bool {
	fmt.Printf("Network.SendMessage: %#v\n", msg)
	return true
}