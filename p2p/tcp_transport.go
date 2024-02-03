package p2p

import (
	"fmt"
	"log"
	"net"
)

type TCPTransportConfig struct {
	ListenAddress    string
	PerformHandshake PerformHandshakeFn  // Function to perform handshake
	Decoder          Decoder             // Decoder to decode incoming messages
	OnPeerConnected  OnPeerConnectedFunc // Function to be called when a new peer is connected

}

type TCPTransport struct {
	config   TCPTransportConfig
	listener net.Listener
	rpcChan  chan RPC
}

func NewTCPTransport(config TCPTransportConfig) *TCPTransport {
	return &TCPTransport{
		config:  config,
		rpcChan: make(chan RPC, 1024),
	}
}

/*
Address implements the Transport interface
returns the address of the transport accepting incoming connections
*/

func (t *TCPTransport) Address() string {
	return t.config.ListenAddress
}

/*
Consume implements the Transport interface
returns a read-only channel of type RPC
which can be used to read incoming messages from another Peer in the network
*/
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcChan
}

/*
Close implements the Transport interface
closes the underlying listener
*/
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

/*
Dial implements the Transport interface
dials the remote node and establishes a TCP connection
*/
func (t *TCPTransport) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}
	t.handleConnection(conn, true)
	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.config.ListenAddress)
	if err != nil {
		return err
	}
	// accept incoming connections
	t.startAcceptingConnectionsLoop()

	log.Printf("TCP transport listening on port: %s\n", t.config.ListenAddress)
	return nil
}

func (t *TCPTransport) startAcceptingConnectionsLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %s", err.Error())
		}
		// handle connection
		go t.handleConnection(conn, false)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()
	// TODO: fix the handshake perform with the remote peer
	if err := t.config.PerformHandshake(conn); err != nil {
		fmt.Printf("Error performing handshake: %s", err.Error())
		conn.Close()
	}
	// handle connection
	peer := NewTCPPeer(conn, outbound)
	fmt.Printf("New Incoming connection from %+v\n", peer)
	if t.config.OnPeerConnected != nil {
		if err = t.config.OnPeerConnected(peer); err != nil {
			return
		}
	}
	// read loop
	rpc := RPC{}
	for {
		if err := t.config.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("Error decoding message: %s", err.Error())
			continue
		}
		rpc.From = conn.RemoteAddr().String()
		if rpc.Stream {
			peer.wg.Add(1)
			fmt.Printf("[%s] incoming stream, waiting...\n", conn.RemoteAddr())
			peer.wg.Wait()
			fmt.Printf("[%s] stream closed, resuming read loop\n", conn.RemoteAddr())
			continue
		}
		t.rpcChan <- rpc
	}
}
