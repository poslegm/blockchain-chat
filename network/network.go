package network

import (
	"fmt"
	"net"
)

type Node struct {
	tcp *net.TCPConn
}

func (n Node) send(msg NetworkMessage) (int, error) {
	return n.tcp.Write(msg)
}

// TODO нужна ли очередь на рассылку?
type NetworkUser struct {
	Nodes 		   map[string]*Node
	Address 	   string
	ConnectQueue	   string
	BroadcastQueue     chan NetworkMessage
}

var CurrentNetworkUser = NetworkUser{}

func setupNetwork() *NetworkUser {
	return &NetworkUser{
		Nodes:map[string]*Node{},
		BroadcastQueue:make(chan NetworkMessage),
		ConnectQueue:make(chan string),
	}
}
// запросы на подключение всегда слушать по tcp, потом отправлять
// в очередь и устанавливать с данными адресами tcp соединение;
// сообщения по tcp читать и писать в базу
func AddNode(node *Node) {
	key := node.tcp.RemoteAddr().String()

	if key != CurrentNetworkUser.Address && node[key] == nil {
		fmt.Println("Node connected", key)
		node[key] = node

		go node.listen()
	}
}

func (node *Node) listen() {
// TODO здесь слушаются сообщения от узла
}

func (networkUser *NetworkUser) listenAll() {
// TODO здесь слушать на подключение
}

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

func Run() {
	CurrentNetworkUser = setupNetwork()
}