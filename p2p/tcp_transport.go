package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPTransportConfig struct {
	ListenAddress    string
	PerformHandshake PerformHandshakeFn // Function to perform handshake
	Decoder          Decoder            // Decoder to decode incoming messages

}

type TCPTransport struct {
	config   TCPTransportConfig
	listener net.Listener
	peerLock sync.RWMutex
	peers    map[net.Addr]Peer
}

func NewTCPTransport(config TCPTransportConfig) *TCPTransport {
	return &TCPTransport{
		config: config,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.config.ListenAddress)
	if err != nil {
		return err
	}
	// accept incoming connections
	t.startAcceptingConnectionsLoop()
	return nil
}

func (t *TCPTransport) startAcceptingConnectionsLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err.Error())
		}
		// handle connection
		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	if err := t.config.PerformHandshake(conn); err != nil {
		fmt.Printf("Error performing handshake: %s", err.Error())
		conn.Close()
	}
	// handle connection
	peer := NewTCPPeer(conn, true)
	fmt.Printf("New Incoming connection from %+v\n", peer)
	// read loop
	msg := &Message{}
	for {
		if err := t.config.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("Error decoding message: %s", err.Error())
			continue
		}
		msg.From = conn.RemoteAddr()
		fmt.Printf("Received message: %+v\n", msg)
	}
}
