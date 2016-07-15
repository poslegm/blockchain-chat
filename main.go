package main

import (
	"fmt"
	"github.com/poslegm/blockchain-chat/db"
	"github.com/poslegm/blockchain-chat/message"
	"github.com/poslegm/blockchain-chat/network"
	"github.com/poslegm/blockchain-chat/server"
	"time"
)

// TODO пока по умолчанию берётся просто первая пара ключей из массива
// TODO получать из базы свой ключ
// TODO сделать интерфейс для записи контактов
// TODO сделать интерфейс для записи своей пары ключей
// TODO сделать нормальное получение ip
// TODO хранить отправленные сообщения в расшифрованном виде

const (
	pubKey     = "message/samplekey_pub.asc"
	privKey    = "message/samplekey_priv.asc"
	passphrase = "sample-key"
)

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Println("main.Run: can't init database ", err.Error())
		return
	}

	kp, err := message.KeyPairFromFile(pubKey, privKey, passphrase)
	if err != nil {
		fmt.Printf("main.Run: cannot create keypair from file: %s\n", err)
		return
	}
	err = db.AddKeys([]*message.KeyPair{kp})
	if err != nil {
		fmt.Println("main.Run: cannot add keypair to db:", err)
		return
	}

	keyPairs, err := db.GetAllKeys()
	if err != nil {
		fmt.Println("main.Run: can't get keys from db ", err.Error())
		return
	} else if len(keyPairs) == 0 {
		fmt.Println("main.Run: there is no key pairs in db")
		return
	}

	err = network.Run(keyPairs[0])
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
			fmt.Println(msg)
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
