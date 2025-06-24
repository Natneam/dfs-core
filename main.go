package main

import (
	"log"

	"natneam.github.io/dfs-core/p2p"
)

func main() {
	tr := p2p.NewTCPTransporter(":3000")

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
