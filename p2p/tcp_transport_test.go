package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTCPTransport(t *testing.T) {
	tcpConfig := TCPTransportConfig{
		ListenAddress:    ":8080",
		PerformHandshake: PerformNoHandshake,
		Decoder:          DefaultDecoder{},
	}
	tcpTransport := NewTCPTransport(tcpConfig)
	assert.NotNil(t, tcpTransport)
	assert.Equal(t, ":8080", tcpTransport.config.ListenAddress)
	// Server
	assert.Nil(t, tcpTransport.ListenAndAccept())
}
