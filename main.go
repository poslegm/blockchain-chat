package main

import (
	"fmt"
	"github.com/poslegm/blockchain-chat/db"
	"github.com/poslegm/blockchain-chat/network"
	"github.com/poslegm/blockchain-chat/server"
	"io/ioutil"
	"strings"
	"time"
)

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Println("main.Run: can't init database ", err.Error())
		return
	}
	// временное решение, возможно, потом добавить возможность добавять адреса через клиент
	err = addIPAddressesToDB()
	if err != nil {
		fmt.Println("main.Run: can't add ip addresses to db ", err.Error())
		return
	}

	keyPairs, err := db.GetAllKeys()
	if err != nil {
		fmt.Println("main.Run: can't get keys from db ", err.Error())
		return
	} else if len(keyPairs) == 0 {
		fmt.Println("main.Run: there is no key pairs in db")
	}

	err = network.Run(keyPairs)
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
		case msg := <-network.CurrentNetworkUser.IncomingMessages:
			fmt.Println("handleNetworkChans: ", msg)
			db.AddMessages([]network.NetworkMessage{msg})
		case address := <-network.CurrentNetworkUser.NewNodes:
			db.AddKnownAddresses([]network.NetAddress{{
				time.Now(),
				address,
				network.TCPPort,
			}})
		case msg := <-network.CurrentNetworkUser.OutgoingMessages:
			fmt.Println(msg)
			db.AddMessages([]network.NetworkMessage{msg})
		}
	}
}

func addIPAddressesToDB() error {
	addresses, err := ioutil.ReadFile("ips.txt")
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return nil
	}

	splitted := strings.Split(string(addresses), "\n")
	networkAddresses := make([]network.NetAddress, len(splitted))

	for i, addr := range splitted {
		networkAddresses[i] = network.NetAddress{Ip: addr, Port: network.TCPPort}
	}
	db.AddKnownAddresses(networkAddresses)

	return nil
}
