package p2p

type Handshake func(Peer) error

func NOPHandShakefunc(Peer) error { return nil }
