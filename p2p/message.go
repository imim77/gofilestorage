package p2p

const (
	IncomingMessage = 0x2
	IncomingStream  = 0x1
)

// RPC holds any arbitrary data that is beign sent
// over each transport between two nodes in the network
type RPC struct {
	Payload []byte
	From    string
	Stream  bool
}
