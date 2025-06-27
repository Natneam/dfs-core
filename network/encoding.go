package network

import (
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(reader io.Reader, msg *Message) error {
	// Pick the first byte to determine the kind of message we're receiving.
	peekBuff := make([]byte, 1)
	if _, err := reader.Read(peekBuff); err != nil {
		return err
	}

	if peekBuff[0] == IncomingStream {
		msg.Stream = true
		return nil
	}

	buffer := make([]byte, 1028)
	n, err := reader.Read(buffer)

	if err != nil {
		return err
	}

	msg.Payload = buffer[:n]
	return nil
}
