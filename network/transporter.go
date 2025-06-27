package network

import (
	"net"
)

// Peer is a representation of a node in the fs network
type Peer interface {
	net.Conn
	RemoteAddr() net.Addr
	Send([]byte) error
	CloseStream()
}

// Transporter handles the communication between nodes in the network.
// This can be TCP, UDP, Websocket or other kind of connections.
type Transporter interface {
	ListenAndAccept() error
	Consume() <-chan Message
	Close() error
	Dial(string) error
}
