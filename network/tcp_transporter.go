package network

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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

func (t *TCPPeer) Close() error {
	return t.conn.Close()
}

func (t *TCPPeer) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

type TCPTransporterOpts struct {
	ListenAddress string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransporter struct {
	TCPTransporterOpts
	listener net.Listener
	msgChan  chan Message
}

func NewTCPTransporter(opts TCPTransporterOpts) *TCPTransporter {
	return &TCPTransporter{
		TCPTransporterOpts: opts,
		msgChan:            make(chan Message),
	}

}

func (t *TCPTransporter) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConn(conn, true)

	return nil
}

// Consume returns read only channel which will be used to read
// messages from another peer node.
func (t *TCPTransporter) Consume() <-chan Message {
	return t.msgChan
}

func (t *TCPTransporter) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}

	go t.ListenAndAcceptLoop()

	log.Printf("TCP Transport listening on port %s", t.ListenAddress)

	return nil
}

func (t *TCPTransporter) Close() error {
	return t.listener.Close()
}

func (t *TCPTransporter) ListenAndAcceptLoop() {

	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP Accept Error : %s", err)
		}

		go t.handleConn(conn, false)
	}

}

func (t *TCPTransporter) handleConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("Error Occurred. Dropping Peer Connection. %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if err = t.OnPeer(peer); err != nil {
		return
	}

	// Read from the peer
	msg := Message{}
	for {
		if err = t.Decoder.Decode(conn, &msg); err != nil {
			// Handle EOF error gracefully
			if err == io.EOF {
				fmt.Printf("EOF : %s\n", err)
				return
			}

			fmt.Printf("TCP Decoding Error : %s\n", err)
			continue
		}

		msg.From = conn.RemoteAddr()

		// Feed the message to chan
		t.msgChan <- msg
	}
}
