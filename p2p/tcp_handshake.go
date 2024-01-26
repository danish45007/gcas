package p2p

import (
	"fmt"
	"net"
)

type TCPHandshake interface {
	PerformHandshakeFn(conn net.Conn) error
}

type PerformHandshakeFn func(conn net.Conn) error

func PerformHandshake(conn net.Conn) error {
	// Send a "Hello" message to the peer
	_, err := conn.Write([]byte("Hello"))
	if err != nil {
		return fmt.Errorf("failed to send Hello message: %s", err.Error())
	}

	// Receive a response from the peer
	buffer := make([]byte, 5) // Adjust the buffer size accordingly
	_, err = conn.Read(buffer)
	if err != nil {
		return fmt.Errorf("failed to receive handshake response: %s", err.Error())
	}

	// Check if the response is as expected
	if string(buffer) != "Ack\n" {
		return fmt.Errorf("unexpected handshake response")
	}

	return nil
}

func PerformNoHandshake(conn net.Conn) error {
	return nil
}
