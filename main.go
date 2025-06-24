package main

import (
	"fmt"
	"log"

	"natneam.github.io/dfs-core/network"
)

func OnPeer(peer network.Peer) error {
	// Do something with the peer
	return nil
}

func main() {
	opts := network.TCPTransporterOpts{
		ListenAddress: ":3000",
		HandshakeFunc: network.NOPHandshakeFunc,
		Decoder:       network.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := network.NewTCPTransporter(opts)

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
