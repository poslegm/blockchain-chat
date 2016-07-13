package network

import (
	"fmt"
	"net"
	"io"
	"encoding/json"
	"os"
)

const TCPPort string = "9005"

type Node struct {
	tcp *net.TCPConn
	key string // TCP address
}

func (n Node) send(msg NetworkMessage) (int, error) {
	marshallMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Network.send: ", err.Error())
		return -1, err
	}
	return n.tcp.Write(marshallMsg)
}

// TODO когда пользователь подключается к сети, его адрес должен быть отправлен всем нодам

type NetworkUser struct {
	Nodes 		   map[string]*Node
	Address 	   string
	ConnectQueue	   chan string
	IncomingMessages   chan NetworkMessage
}

var CurrentNetworkUser *NetworkUser = new(NetworkUser)
// создаёт объект для обработчика сети
func setupNetwork(address string) *NetworkUser {
	networkUser := new(NetworkUser)
	networkUser.Nodes = map[string]*Node{}
	networkUser.Address = address
	networkUser.ConnectQueue = make(chan string)
	networkUser.IncomingMessages = make(chan NetworkMessage)
	return networkUser
}

// добавляет узел и запускает горутину на его прослушивание
func (networkUser *NetworkUser) addNode(node *Node) {
	if node.key != networkUser.Address && networkUser.Nodes[node.key] == nil {
		fmt.Println("Node connected", node.key)
		networkUser.Nodes[node.key] = node

		go node.listen(networkUser)
	}
}

func (networkUser *NetworkUser) removeNode(node *Node) {
	node.tcp.Close()
	delete(networkUser.Nodes, node.key)
}

func currentUserRemoveNode(node *Node) {
	CurrentNetworkUser.removeNode(node)
}

// слушает сообщения от указанного узла и пишет их в базу
func (node *Node) listen(networkUser *NetworkUser) {
	for {
		buffer, err := node.listenTCP()

		if err == io.EOF {
			break
		}

		msg := new(NetworkMessage)
		err = json.Unmarshal(buffer, msg)
		if err != nil {
			fmt.Println("Network.node.listen: ", err.Error())
			continue
		}

		networkUser.IncomingMessages <- *msg
	}
}

// запись одного сообщения по TCP в буфер (вспомогательная процедура)
func (node *Node) listenTCP() ([]byte, error) {
	buffer := make([]byte, 4096)
	tmp := make([]byte, 256)

	for {
		n, err := node.tcp.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Network.node.listen: ", err.Error())
				currentUserRemoveNode(node)
			}
			return nil, err
		}
		buffer = append(buffer, tmp[:n]...)
	}

	return buffer, nil
}

// прослушивает tcp на предмет запросов на подключение
func (networkUser *NetworkUser) listenTCPRequests() {
	address, err := net.ResolveTCPAddr("tcp4", networkUser.Address)
	if err != nil {
		fmt.Println("Network.listenTCPRequests: ", err.Error())
		return
	}

	listener, err := net.ListenTCP("tcp4", address)
	if err != nil {
		fmt.Println("Network.listenTCPRequests: ", err.Error())
		return
	}

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Network.listenTCPRequests: ", err.Error())
			return
		}

		networkUser.addNode(&Node{connection, connection.RemoteAddr().String()})
	}
}

// отправка запросов на соединение с узлами из очереди подключений
func (networkUser *NetworkUser) sendConnectionRequests() {
	for {
		address := <- networkUser.ConnectQueue
		address = address + ":" + TCPPort

		if address != networkUser.Address && networkUser.Nodes[address] == nil {
			go networkUser.setConnection(address)
		}
	}
}

// установка соединения с узлом по адресу
func (networkUser *NetworkUser) setConnection(address string) {
	tcpAddress, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		fmt.Println("Network.setConnection: ", err.Error())
		return
	}
	// TODO возможно, заморочиться с таймаутами на соединение
	connection, err := net.DialTCP("tcp", nil, tcpAddress)
	if err != nil {
		fmt.Println("Network.setConnection: ", err.Error())
		return
	}

	fmt.Println("Network.setConnection: connected ", connection.RemoteAddr())
	networkUser.addNode(&Node{connection, connection.RemoteAddr().String()})
}

// рассылка сообщения по узлам
func (network *NetworkUser) SendMessage(msg NetworkMessage) {
	fmt.Printf("Network.SendMessage: %#v\n", msg)

	for k, node := range network.Nodes {
		fmt.Println("Broadcasting...", k)

		go func() {
			_, err := node.send(msg)
			if err != nil {
				fmt.Println("Error broadcasting to ", node.tcp.RemoteAddr())
			}
		}()
	}
}

func currentAddress() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	address, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}

	fmt.Println("Network.currentAddress: address ", address)
	return address[0] + ":" + TCPPort, nil
}

func Run() error {
	address, err := currentAddress()
	if err != nil {
		fmt.Println("Network.Run: ", err.Error())
		fmt.Println("Can't create network")
		return err
	}

	CurrentNetworkUser = setupNetwork(address)
	go CurrentNetworkUser.sendConnectionRequests()
	go CurrentNetworkUser.listenTCPRequests()

	return nil
}