package p2p

import (
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(reader io.Reader, msg *Message) error {
	buffer := make([]byte, 1028)
	n, err := reader.Read(buffer)

	if err != nil {
		return err
	}

	msg.Payload = buffer[:n]
	return nil
}
