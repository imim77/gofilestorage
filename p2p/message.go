package p2p

import "net"

// Message holds any arbitrary data that is beign sent
// over each transport between two nodes in the network
type Message struct {
	Payload []byte
	From    net.Addr
}
