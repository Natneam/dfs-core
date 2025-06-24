package network

type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(p Peer) error { return nil }
