package main

import (
	"fmt"
	"log"
	"time"

	"natneam.github.io/dfs-core/cipher"
	"natneam.github.io/dfs-core/cli"
	"natneam.github.io/dfs-core/network"
	"natneam.github.io/dfs-core/server"
	"natneam.github.io/dfs-core/store"
)

func makeFileServer(addr string, nodes ...string) *server.FileServer {
	tcpTransporterOpts := network.TCPTransporterOpts{
		ListenAddress: addr,
		HandshakeFunc: network.NOPHandshakeFunc,
		Decoder:       network.DefaultDecoder{},
	}

	tcpTransporter := network.NewTCPTransporter(tcpTransporterOpts)

	fileServerOpts := server.FileServerOpts{
		StorageRoot:       addr + "_files",
		PathTransformFunc: store.HashPathTransformFunc,
		Transporter:       tcpTransporter,
		BootstrapNodes:    nodes,
		EncKey:            cipher.NewEncryptionKey(),
	}

	s := server.NewFileServer(fileServerOpts)
	tcpTransporter.OnPeer = s.OnPeer
	return s

}

func main() {
	port, nodes, err := cli.Start()

	if err != nil {
		log.Fatal(err)
	}

	fs := makeFileServer(fmt.Sprintf(":%d", port), nodes...)

	go func() {
		fs.Start()
	}()

	time.Sleep(time.Second) // Wait for the server to start.

	cli.InteractiveCli(fs)
}
