package main

import (
	"github.com/poslegm/blockchain-chat/server"
	"github.com/poslegm/blockchain-chat/db"
	"github.com/poslegm/blockchain-chat/network"
	"fmt"
	"time"
)
// TODO запись публичных ключей в базе
// TODO полноценный обмен данными с клиентом

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Println("main.Run: can't init database ", err.Error())
		return
	}

	err = network.Run()
	if err != nil {
		fmt.Println("main.Run: can't run network ", err.Error())
		return
	}

	go createConnectQueue()
	go handleNetworkChans()

	server.Run("./client", "8080")
}

func createConnectQueue() {
	knownAdresses, err := db.GetKnownAddresses()
	if err != nil {
		fmt.Println("main.Run: can't get addresses from db ", err.Error())
	}

	for _, address := range knownAdresses {
		network.CurrentNetworkUser.ConnectQueue <- address.Ip
	}
}

func handleNetworkChans() {
	for {
		select {
		case msg := <- network.CurrentNetworkUser.IncomingMessages:
			db.AddMessages([]network.NetworkMessage{msg})
		case address := <- network.CurrentNetworkUser.NewNodes:
			db.AddKnownAddresses([]network.NetAddress{{
				time.Now(),
				address,
				network.TCPPort,
			}})
		}
	}
}