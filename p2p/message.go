package p2p

import "net"

// Message represent the any arbitrary data sent b/w nodes in the network
type Message struct {
	Payload []byte
	From    net.Addr
}
