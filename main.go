package main

import (
	"fmt"
	"log"

	"natneam.github.io/dfs-core/p2p"
)

func OnPeer(peer p2p.Peer) error {
	// Do something with the peer
	return nil
}

func main() {
	opts := p2p.TCPTransporterOpts{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransporter(opts)

	go func() {
		for {
			fmt.Printf("Message : %+v\n", <-tr.Consume())
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
