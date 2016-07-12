package network

import (
	"time"
	"net"
)

type NetAddress struct {
	//last online time of peer
	lastseen time.Time;

	//peer's ip address
	ip net.IP

	//peer's port
	port uint16
}
