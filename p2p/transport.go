package p2p

// Peer represent the remote node in the network
type Peer interface {
	// net.Conn
	CloseStream()
	Send([]byte) error
}

// Transport is responsible for handling the communication b/w nodes in the network
// Transport type can be implemented using (TCP or UDP or Websocket...)
type Transport interface {
	Address() string
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
