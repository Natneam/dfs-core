package p2p

type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(p Peer) error { return nil }
