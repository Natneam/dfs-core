package main

import (
	"fmt"
	"io"
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
	s2 := makeFileServer(":11000")
	s3 := makeFileServer(":12000", ":11000", ":10000")

	go func() { log.Fatal(s1.Start()) }()
	time.Sleep(500 * time.Millisecond)
	go func() { log.Fatal(s2.Start()) }()

	time.Sleep(2 * time.Second)
	go func() { log.Fatal(s3.Start()) }()
	time.Sleep(2 * time.Second)

	// for i := range 10 {
	// 	data, _ := os.Open("./main.go")
	// 	s3.Store(fmt.Sprintf("%s_%d", "mykey", i), data)
	// 	data.Close()
	// 	time.Sleep(time.Microsecond * 100)
	// }

	_, r, err := s3.Get("mykey_0")
	if err != nil {
		fmt.Printf("data not found : %+v\n", err)
	}
	b, _ := io.ReadAll(r)
	println(string(b))
}
