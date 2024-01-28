package cas

import (
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
	config   FileServerConfig
	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	storage  *storage.Storage
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
		config:  config,
		storage: storage.NewStorage(storeConfig),
		peers:   make(map[string]p2p.Peer),
		channel: make(chan struct{}),
	}
}
