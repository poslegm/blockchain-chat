package network

import (
	"encoding/json"
	"fmt"
	"github.com/poslegm/blockchain-chat/message"
	"io"
	"net"
	"os"
)

const TCPPort string = "9005"

type Node struct {
	tcp *net.TCPConn
	key string // TCP address
}

func (n Node) send(msg NetworkMessage) error {
	encoder := json.NewEncoder(n.tcp)

	fmt.Println("Sended: ", msg)
	return encoder.Encode(msg)
}

type NetworkUser struct {
	Nodes            map[string]*Node // контакты
	Address          string
	ConnectQueue     chan string         // очередь на отправку запросов соединения
	IncomingMessages chan NetworkMessage // входящие для добавления в базу
	OutgoingMessages chan NetworkMessage
	NewNodes         chan string // адреса новых соединений для добавления в базу
	KeyPairs         []*message.KeyPair
}

var CurrentNetworkUser *NetworkUser = new(NetworkUser)

// создаёт объект для обработчика сети
func setupNetwork(address string, kps []*message.KeyPair) *NetworkUser {
	networkUser := new(NetworkUser)
	networkUser.Nodes = map[string]*Node{}
	networkUser.Address = address
	networkUser.ConnectQueue = make(chan string)
	networkUser.IncomingMessages = make(chan NetworkMessage)
	networkUser.OutgoingMessages = make(chan NetworkMessage)
	networkUser.KeyPairs = kps
	return networkUser
}

// добавляет узел и запускает горутину на его прослушивание
func (networkUser *NetworkUser) addNode(node *Node) {
	if node.key != networkUser.Address && networkUser.Nodes[node.key] == nil {
		fmt.Println("Node connected ", node.key)
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
		message, err := node.listenTCP()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("network.listen: can't receive message ", err.Error())
			return
		}

		fmt.Println("network.listen: received ", *message)

		switch message.MessageType {
		case MESSAGE:
			fmt.Println("WRITED: ", *message)
			networkUser.IncomingMessages <- *message
		case REQUEST:
			fmt.Println("REQUEST", *message)
			networkUser.ConnectQueue <- message.IP
		}
	}
}

// запись одного сообщения по TCP в буфер (вспомогательная процедура)
func (node *Node) listenTCP() (*NetworkMessage, error) {
	decoder := json.NewDecoder(node.tcp)

	msg := new(NetworkMessage)
	err := decoder.Decode(msg)
	if err != nil {
		fmt.Println("network.listenTCP: can't decode message ", err.Error())
		fmt.Println(msg)
		currentUserRemoveNode(node)
	}
	return msg, err
}

// прослушивает tcp на предмет запросов на подключение
func (networkUser *NetworkUser) listenTCPRequests() {
	address, err := net.ResolveTCPAddr("tcp4", networkUser.Address)
	if err != nil {
		fmt.Println("Network.listenTCPRequests: can't resolve addr", err.Error())
		return
	}

	listener, err := net.ListenTCP("tcp4", address)
	if err != nil {
		fmt.Println("Network.listenTCPRequests: can't listen tcp ", err.Error())
		return
	}

	for {
		connection, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Network.listenTCPRequests: can't accept tcp ", err.Error())
			return
		}

		networkUser.addNode(&Node{connection, connection.RemoteAddr().String()})
		networkUser.NewNodes <- connection.RemoteAddr().String()

		msg := NetworkMessage{
			MessageType: REQUEST,
			IP:          connection.RemoteAddr().String(),
		}

		networkUser.OutgoingMessages <- msg
		go networkUser.SendMessage(msg)
	}
}

// отправка запросов на соединение с узлами из очереди подключений
func (networkUser *NetworkUser) sendConnectionRequests() {
	for {
		address := <-networkUser.ConnectQueue
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

	for k, node := range network.Nodes {
		fmt.Println("Broadcasting...", k)

		go func() {
			err := node.send(msg)
			if err != nil {
				fmt.Println("Error broadcasting ("+err.Error()+") to ", node.tcp.RemoteAddr())
			}
		}()
	}
}

// получает текущий локальный ip
func currentAddress() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	address, err := net.LookupHost(name)
	if err != nil {
		return "", err
	}

	addrId := 0
	if len(address) != 1 {
		for i, addr := range address {
			binAddr := net.ParseIP(addr).To4()

			if binAddr != nil && len(binAddr) == 4 {
				addrId = i
				break
			}
		}
	}
	fmt.Println("Network.currentAddress: address ", address[addrId]+":"+TCPPort)
	return address[addrId] + ":" + TCPPort, nil
}

func Run(kps []*message.KeyPair) error {
	address, err := currentAddress()
	if err != nil {
		fmt.Println("Network.Run: can't get current IP address ", err.Error())
		fmt.Println("Can't create network")
		return err
	}

	CurrentNetworkUser = setupNetwork(address, kps)
	go CurrentNetworkUser.sendConnectionRequests()
	go CurrentNetworkUser.listenTCPRequests()

	return nil
}
