package main

import (
	"bytes"
	"fmt"
	"log"
	"time"

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
	}

	s := server.NewFileServer(fileServerOpts)
	tcpTransporter.OnPeer = s.OnPeer
	return s

}

func main() {
	s1 := makeFileServer(":10000")
	s2 := makeFileServer(":11000", ":10000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second * 2)
	go s2.Start()
	time.Sleep(time.Second * 2)

	for i := range 10 {
		data := bytes.NewReader([]byte("Hello this is a large data"))
		s2.Store(fmt.Sprintf("%s_%d", "mykey", i), data)
		time.Sleep(time.Microsecond * 100)
	}

	// r, err := s2.Get("mykey")
	// if err != nil {
	// 	log.Fatal("data not found")
	// }
	// b, _ := io.ReadAll(r)
	// println(string(b))
	select {}
}
