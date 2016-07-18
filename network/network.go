package network

import (
	"encoding/json"
	"fmt"
	"github.com/poslegm/blockchain-chat/message"
	"io"
	"net"
	"net/http"
	"os"
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
	marshallMsg = append(marshallMsg, 4)
	fmt.Println(marshallMsg)
	return n.tcp.Write(marshallMsg)
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
		buffers, err := node.listenTCP()

		if err == io.EOF {
			break
		}

		messages := make([]*NetworkMessage, 0)
		for _, buffer := range buffers {
			msg := new(NetworkMessage)
			err = json.Unmarshal(buffer, msg)
			if err != nil {
				fmt.Println("Network.node.listen: ", err.Error())
				continue
			}
			messages = append(messages, msg)
		}

		for _, msg := range messages {
			switch msg.MessageType {
			case MESSAGE:
				fmt.Println("WRITED: ", *msg)
				networkUser.IncomingMessages <- *msg
			case REQUEST:
				fmt.Println("REQUEST", msg)
				networkUser.ConnectQueue <- msg.IP
			}
		}
	}
}

// запись одного сообщения по TCP в буфер (вспомогательная процедура)
func (node *Node) listenTCP() ([][]byte, error) {
	buffer := make([]byte, 4096)

	_, err := node.tcp.Read(buffer)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Network.node.listen: ", err.Error())
			currentUserRemoveNode(node)
		}
		return nil, err
	}

	messages := make([][]byte, 0)
	sep := 0
	for i, c := range buffer {
		if c == 4 {
			messages = append(messages, buffer[sep:i])
			sep = i + 1
		}
	}
	return messages, nil
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

		//if networkUser.Nodes[connection.RemoteAddr().String()]
		networkUser.addNode(&Node{connection, connection.RemoteAddr().String()})
		networkUser.NewNodes <- connection.RemoteAddr().String()
		// TODO возможно, по сети бесконечно будет летать это сообщение
		go networkUser.SendMessage(NetworkMessage{
			MessageType: REQUEST,
			IP:          connection.RemoteAddr().String(),
		})
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
	fmt.Println("Network.currentAddress: address ", address[addrId])
	return address[addrId] + ":" + TCPPort, nil
}

// получает текущий глобальный ip, используя сервис Амазона
func currentGlobalAddress() (string, error) {
	resp, err := http.Get("http://checkip.amazonaws.com")
	if err != nil {
		return "", err
	}

	address := make([]byte, 20)
	n, err := resp.Body.Read(address)
	if err != nil && err != io.EOF {
		return "", err
	}
	address = address[:n-1] // n - 1, потому что на конце стоит символ переноса строки

	addressString := string(address)
	addressString += ":" + TCPPort
	fmt.Println("Network.currentAddress: address ", addressString, n)

	return addressString, nil
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
