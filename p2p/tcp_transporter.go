package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn

	// outbound = false, if we're accept a connection request
	// outbound = true, if we're dialing a connection request
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransporter struct {
	listenAddress string
	listener      net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransporter(listenAdd string) *TCPTransporter {
	return &TCPTransporter{
		listenAddress: listenAdd,
	}

}

func (t *TCPTransporter) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}

	go t.ListenAndAcceptLoop()

	return nil
}

func (t *TCPTransporter) ListenAndAcceptLoop() {

	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Accept Error : %s", err)
		}

		go t.handleConn(conn)
	}

}

func (t *TCPTransporter) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)
	fmt.Printf("New Incoming connection %+v \n", peer)
}
