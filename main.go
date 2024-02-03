package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"time"

	cas "github.com/danish45007/gcas/cas"
	"github.com/danish45007/gcas/cascrypto"
	"github.com/danish45007/gcas/p2p"
	"github.com/danish45007/gcas/storage"
)

func makeServer(listenAddr string, nodes ...string) *cas.FileServer {
	tcpConfig := p2p.TCPTransportConfig{
		ListenAddress:    listenAddr,
		PerformHandshake: p2p.PerformNoHandshake,
		Decoder:          p2p.DefaultDecoder{},
	}
	fileServerConfig := cas.FileServerConfig{
		EncryptionKey:  cascrypto.NewEncryptionKey(),
		StorageRoot:    listenAddr + "_newtwork",
		PathTransform:  storage.CASPathTransform,
		BootstrapNodes: nodes,
	}
	fs := cas.NewFileServer(fileServerConfig)
	tcpConfig.OnPeerConnected = fs.OnPeerConnected
	return fs
}

func main() {
	// node1
	fs1 := makeServer(":3000", "")
	// node2
	fs2 := makeServer(":7000", "")
	// node3
	fs3 := makeServer(":8000", ":3000", ":7000")
	// start the server 1
	go func() {
		log.Fatal(fs1.Start())
	}()
	// add delay
	time.Sleep(500 * time.Millisecond)
	// start the server 2
	go func() {
		log.Fatal(fs2.Start())
	}()
	// add delay
	time.Sleep(2 * time.Second)
	// start the server 3
	go func() {
		log.Fatal(fs3.Start())
	}()
	// add delay
	time.Sleep(2 * time.Second)
	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("test_%d.png", i)
		data := bytes.NewReader([]byte("test data"))
		fs3.StoreFile(key, data)
		if err := fs3.Storage.Delete(fs3.Config.ID, key); err != nil {
			log.Fatal(err)
		}

		r, err := fs3.GetFile(key)
		if err != nil {
			log.Fatal(err)
		}

		b, err := io.ReadAll(r)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(b))
	}
}
