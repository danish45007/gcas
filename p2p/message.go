package p2p

const (
	IncomingMessage = 0x01
	OutgoingMessage = 0x02
)

// RPC represent the any arbitrary data sent b/w nodes in the network
type RPC struct {
	Payload []byte
	From    string
	Stream  bool
}
