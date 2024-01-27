package p2p

import "net"

// RPC represent the any arbitrary data sent b/w nodes in the network
type RPC struct {
	Payload []byte
	From    net.Addr
}
