package p2p

// Peer represent the remote node in the network
type Peer interface {
	Close() error
}

// Transport is responsible for handling the communication b/w nodes in the network
// Transport type can be implemented using (TCP or UDP or Websocket...)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
