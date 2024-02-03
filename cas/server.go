package cas

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/danish45007/gcas/cascrypto"
	"github.com/danish45007/gcas/p2p"
	"github.com/danish45007/gcas/storage"
)

type FileServerConfig struct {
	ID             string
	EncryptionKey  []byte
	StorageRoot    string
	PathTransform  storage.PathTransformFunc
	Transport      p2p.Transport
	BootstrapNodes []string
}

type FileServer struct {
	Config   FileServerConfig
	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	Storage  *storage.Storage
	channel  chan struct{}
}

func NewFileServer(config FileServerConfig) *FileServer {
	storeConfig := storage.StorageConfig{
		PathTransform: config.PathTransform,
		SetRoot:       config.StorageRoot,
	}
	if len(config.ID) == 0 {
		config.ID = cascrypto.GenerateId()
	}
	return &FileServer{
		Config:  config,
		Storage: storage.NewStorage(storeConfig),
		peers:   make(map[string]p2p.Peer),
		channel: make(chan struct{}),
	}
}

func (fs *FileServer) Stop() {
	close(fs.channel)
}

func (fs *FileServer) OnPeerConnected(peer p2p.Peer) error {
	// acquire the lock
	fs.peerLock.Lock()
	defer fs.peerLock.Unlock()
	// add the peer to the map
	fs.peers[peer.Address()] = peer
	log.Printf("[%s] Peer connected: %s\n", fs.Config.Transport.Address(), peer.Address())
	return nil
}

func (fs *FileServer) loop() {
	defer func() {
		log.Println("File server loop stopped due to an error or quit")
		fs.Config.Transport.Close()
	}()
	for {
		select {
		case <-fs.channel:
			return
		case rpc := <-fs.Config.Transport.Consume():
			var msg BroadCastMessage
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("decoding error: ", err)
			}
			if err := fs.handleMessage(rpc.From, &msg); err != nil {
				log.Println("handle message error: ", err)
			}

		}
	}
}

func (fs *FileServer) handleMessage(from string, msg *BroadCastMessage) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		// store the file in the local storage
		fs.handleMessageStoreFile(from, v)
	case MessageGetFile:
		// get the file from the local storage
		fs.handleMessageGetFile(from, v)
	}
	return nil
}

func (fs *FileServer) handleMessageGetFile(from string, msg MessageGetFile) error {
	if !fs.Storage.Has(msg.ID, msg.Key) {
		return fmt.Errorf("[%s] need to serve file (%s) but it does not exist on disk", fs.Config.Transport.Address(), msg.Key)
	}

	fmt.Printf("[%s] serving file (%s) over the network\n", fs.Config.Transport.Address(), msg.Key)

	fileSize, r, err := fs.Storage.Read(msg.ID, msg.Key)
	if err != nil {
		return err
	}

	if rc, ok := r.(io.ReadCloser); ok {
		fmt.Println("closing readCloser")
		defer rc.Close()
	}

	peer, ok := fs.peers[from]
	if !ok {
		return fmt.Errorf("peer %s not in map", from)
	}

	// First send the "incomingStream" byte to the peer and then we can send
	// the file size as an int64.
	peer.Send([]byte{p2p.IncomingMessage})
	binary.Write(peer, binary.LittleEndian, fileSize)
	n, err := io.Copy(peer, r)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] written (%d) bytes over the network to %s\n", fs.Config.Transport.Address(), n, from)

	return nil
}

func (fs *FileServer) handleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := fs.peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) could not be found in the peer list", from)
	}

	n, err := fs.Storage.Write(msg.ID, msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}

	fmt.Printf("[%s] written %d bytes to disk\n", fs.Config.Transport.Address(), n)

	peer.CloseStream()

	return nil
}

func (fs *FileServer) bootstrapNetwork() error {
	for _, addr := range fs.Config.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		go func(addr string) {
			fmt.Printf("[%s] attempting to connect with remote %s\n", fs.Config.Transport.Address(), addr)
			if err := fs.Config.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)
			}
		}(addr)
	}

	return nil
}

func (fs *FileServer) Start() error {
	fmt.Printf("[%s] starting fileserver...\n", fs.Config.Transport.Address())

	if err := fs.Config.Transport.ListenAndAccept(); err != nil {
		return err
	}

	fs.bootstrapNetwork()

	fs.loop()

	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
	gob.Register(MessageGetFile{})
}
