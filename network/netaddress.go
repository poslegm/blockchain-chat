package network

import (
	"time"
)

type NetAddress struct {
	//last online time of peer
	Lastseen time.Time;

	//peer's ip address
	Ip string

	//peer's port
	Port string
}
