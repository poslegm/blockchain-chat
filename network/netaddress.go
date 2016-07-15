package network

import (
	"time"
)

type NetAddress struct {
	//last online time of peer
	Lastseen time.Time;

	//peer's ip address
	Ip       string

	//peer's port
	Port     string
}

func CreateNetAddress(lastseen time.Time, ip, port string) NetAddress {
	return NetAddress{Lastseen:lastseen, Ip:ip, Port:port}
}