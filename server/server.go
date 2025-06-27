package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
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
	gob.Register(&GetMessagePayload{})
	gob.Register(&StoreMessagePayload{})

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

	if err := s.bootstrapNetwork(); err != nil {
		return err
	}

	s.loop()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitchan)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	fileBuf := new(bytes.Buffer)
	tee := io.TeeReader(r, fileBuf)

	n, err := s.Store.Write(key, tee)

	if err != nil {
		return err
	}

	// Send the Store message with size and key
	msg := DataMessage{
		Payload: StoreMessagePayload{
			Key:  key,
			Size: n,
		},
	}

	msgBuf := new(bytes.Buffer)
	if err := gob.NewEncoder(msgBuf).Encode(msg); err != nil {
		return err
	}

	peers := []io.Writer{}

	for _, peer := range s.peers {
		peer.Send(msgBuf.Bytes())
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	if _, err = io.Copy(mw, fileBuf); err != nil {
		return fmt.Errorf("failed to send file content to peers: %w", err)
	}

	return nil
}

func (s *FileServer) OnPeer(p network.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p

	log.Printf("Connected with remote %s", p.LocalAddr())

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		fmt.Printf("File Server Stopped Due to User Quit Action.\n")
		s.Transporter.Close()
	}()

	for {
		select {
		case rpc := <-s.Transporter.Consume():
			var msg DataMessage
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Fatal(err)
			}
			s.handleMessage(rpc.From.String(), &msg)
		case <-s.quitchan:
			return
		}
	}
}

func (s *FileServer) handleMessage(from string, msg *DataMessage) error {
	switch v := msg.Payload.(type) {
	case *StoreMessagePayload:
		fmt.Printf("Store message received %+v\n", v)

		peer, ok := s.peers[from]
		if !ok {
			return fmt.Errorf("peer (%s) could not be found in the peer list", from)
		}

		_, err := s.Store.Write(v.Key, io.LimitReader(peer, v.Size))
		if err != nil {
			return err
		}
		peer.CloseStream()
		fmt.Printf("Data received and stored to disk %+v\n", v)
	}
	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	if len(s.BootstrapNodes) == 0 {
		return nil
	}
	wg := sync.WaitGroup{}
	for _, node := range s.BootstrapNodes {
		wg.Add(1)
		go func(node string) {
			fmt.Println("Attempting to connect with remote => ", node)
			if err := s.Transporter.Dial(node); err != nil {
				log.Printf("dial error : %s\n", err)
				wg.Done()
				return
			}
			wg.Done()
		}(node)
	}
	wg.Wait()
	return nil
}

type StoreMessagePayload struct {
	Key  string
	Size int64
}

type GetMessagePayload struct {
	Key string
}

type DataMessage struct {
	Payload any
}
