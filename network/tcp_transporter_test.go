package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPtransporter(t *testing.T) {
	opts := TCPTransporterOpts{
		ListenAddress: ":3000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}

	tr := NewTCPTransporter(opts)

	assert.Equal(t, opts.ListenAddress, tr.ListenAddress)
}
