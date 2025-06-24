package p2p

import (
	"fmt"
	"io"
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

type TCPTransporterOpts struct {
	ListenAddress string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransporter struct {
	TCPTransporterOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransporter(opts TCPTransporterOpts) *TCPTransporter {
	return &TCPTransporter{
		TCPTransporterOpts: opts,
	}

}

func (t *TCPTransporter) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
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

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP Handshake Error %s\n", &err)
		return
	}

	// Read from the peer
	msg := &Message{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			// Handle EOF error gracefully
			if err == io.EOF {
				fmt.Printf("EOF : %s\n", err)
				break
			}

			fmt.Printf("TCP Decoding Error : %s\n", err)
			continue
		}

		fmt.Printf("Message : %s\n", msg)
	}
}
