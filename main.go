package main

import (
	"log"

	"natneam.github.io/dfs-core/p2p"
)

func main() {
	opts := p2p.TCPTransporterOpts{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransporter(opts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
