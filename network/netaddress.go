package network

import (
	"time"
	"net"
)

type NetAddress struct {
	//last online time of peer
	Lastseen time.Time;

	//peer's ip address
	Ip net.IP

	//peer's port
	Port uint16
}
