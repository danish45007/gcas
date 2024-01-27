package main

import (
	"fmt"

	"github.com/danish45007/gcas/p2p"
)

func main() {
	tcpConfig := p2p.TCPTransportConfig{
		ListenAddress:    ":8080",
		PerformHandshake: p2p.PerformNoHandshake,
		Decoder:          p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpConfig)
	go func() {
		for {
			rpc := <-tcpTransport.Consume()
			fmt.Printf("Received RPC %+v\n", rpc)
		}
	}()
	if err := tcpTransport.ListenAndAccept(); err != nil {
		fmt.Printf("Error listening and accepting connections: %s", err.Error())
	}
	select {}

}
