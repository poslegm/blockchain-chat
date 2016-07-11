package db

import "github.com/poslegm/blockchain-chat/network"

func GetAllMessages() []network.NetworkMessage {
	return []network.NetworkMessage{
		network.NetworkMessage{"Букер", "VS94SKI", "KEK"},
		network.NetworkMessage{"VS94SKI", "Букер", "CHPEK"},
	}
}