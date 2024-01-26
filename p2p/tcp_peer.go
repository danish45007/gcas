package p2p

import "net"

// TCPPeer represent the remote node over an established TCP connection
type TCPPeer struct {
	conn     net.Conn // conn is the underlying TCP connection of the peer
	outbound bool     // outbound: true if we dial the connection
	// outbound: false if we accept the connection
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}
