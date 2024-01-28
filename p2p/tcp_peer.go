package p2p

import (
	"net"
	"sync"
)

type OnPeerConnectedFunc func(Peer) error

// TCPPeer represent the remote node over an established TCP connection
type TCPPeer struct {
	Conn     net.Conn // conn is the underlying TCP connection of the peer
	outbound bool     // outbound: true if we dial the connection
	// outbound: false if we accept the connection
	wg *sync.WaitGroup
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		wg:       &sync.WaitGroup{},
	}
}

/*
CloseStream implements the Peer interface
responsible for decrementing the waitgroup counter when the stream is closed
*/
func (p *TCPPeer) CloseStream() {
	p.wg.Done()
}

/*
Send implements the Peer interface
responsible for sending the data over the TCP connection
*/
func (p *TCPPeer) Send(data []byte) error {
	_, err := p.Conn.Write(data)
	return err
}
