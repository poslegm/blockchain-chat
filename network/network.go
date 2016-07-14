package network

import (
	"fmt"
	"net"
	"io"
	"encoding/json"
	"os"
	"github.com/poslegm/blockchain-chat/message"
)

const TCPPort string = "9005"

type Node struct {
	tcp *net.TCPConn
	key string // TCP address
}
// TODO покрыть тестами

func (n Node) send(msg NetworkMessage) (int, error) {
	marshallMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("Network.send: ", err.Error())
		return -1, err
	}
	return n.tcp.Write(marshallMsg)
}

type NetworkUser struct {
	Nodes 		   map[string]*Node // контакты
	Address 	   string
	ConnectQueue	   chan string // очередь на отправку запросов соединения
	IncomingMessages   chan NetworkMessage // входящие для добавления в базу
	OutgoingMessages   chan NetworkMessage
	NewNodes           chan string // адреса новых соединений для добавления в базу
	KeyPair		   *message.KeyPair
}

var CurrentNetworkUser *NetworkUser = new(NetworkUser)
// создаёт объект для обработчика сети
func setupNetwork(address string, kp *message.KeyPair) *NetworkUser {
	networkUser := new(NetworkUser)
	networkUser.Nodes = map[string]*Node{}
	networkUser.Address = address
	networkUser.ConnectQueue = make(chan string)
	networkUser.IncomingMessages = make(chan NetworkMessage)
	networkUser.OutgoingMessages = make(chan NetworkMessage)
	networkUser.KeyPair = kp
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

		switch msg.MessageType {
		case MESSAGE: networkUser.IncomingMessages <- *msg
		case REQUEST: networkUser.ConnectQueue <- msg.IP
		}

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
		networkUser.NewNodes <- connection.RemoteAddr().String()
		go networkUser.SendMessage(NetworkMessage{
			MessageType:REQUEST,
			IP:connection.RemoteAddr().String(),
		})
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

	network.OutgoingMessages <- msg

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

	addrId := 0
	if len(address) != 1 {
		for i, addr := range address {
			binAddr := net.ParseIP(addr).To4()

			if(binAddr != nil && len(binAddr) == 4) {
				addrId = i
				break
			}
		}
	}
	fmt.Println("Network.currentAddress: address ", address[addrId])
	return address[addrId] + ":" + TCPPort, nil
}

func Run(kp *message.KeyPair) error {
	address, err := currentAddress()
	if err != nil {
		fmt.Println("Network.Run: ", err.Error())
		fmt.Println("Can't create network")
		return err
	}

	CurrentNetworkUser = setupNetwork(address, kp)
	go CurrentNetworkUser.sendConnectionRequests()
	go CurrentNetworkUser.listenTCPRequests()

	return nil
}