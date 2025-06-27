package network

import "net"

const (
	IncomingMessage = 0x0
	IncomingStream  = 0x1
)

type Message struct {
	From    net.Addr
	Payload []byte
	Stream  bool
}

type DataMessage struct {
	Payload any
}

type StoreMessagePayload struct {
	Key  string
	Size int64
}

type GetMessagePayload struct {
	Key string
}
