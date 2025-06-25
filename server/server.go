package server

import (
	"fmt"
	"log"
	"sync"

	"natneam.github.io/dfs-core/network"
	"natneam.github.io/dfs-core/store"
)

type FileServerOpts struct {
	StorageRoot       string
	Transporter       network.Transporter
	PathTransformFunc store.PathTransformFunc
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]network.Peer

	Store    *store.Store
	quitchan chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := store.StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		peers:          make(map[string]network.Peer),
		Store:          store.NewStore(storeOpts),
		quitchan:       make(chan struct{}),
	}
}

func (s *FileServer) Start() error {
	if err := s.Transporter.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()
	s.loop()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitchan)
}

func (s *FileServer) OnPeer(p network.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p

	log.Printf("Connected with remote %s", p.RemoteAddr())

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		fmt.Printf("File Server Stopped Due to User Quit Action.\n")
		s.Transporter.Close()
	}()

	for {
		select {
		case msg := <-s.Transporter.Consume():
			fmt.Println(msg)
		case <-s.quitchan:
			return
		}
	}
}

func (s *FileServer) bootstrapNetwork() {
	for _, node := range s.BootstrapNodes {
		go func(node string) {
			fmt.Println("Attempting to connect with remote => ", node)
			if err := s.Transporter.Dial(node); err != nil {
				log.Printf("dial error : %s\n", err)
			}
		}(node)
	}
}
