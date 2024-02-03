package cas

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"time"

	"github.com/danish45007/gcas/cascrypto"
	"github.com/danish45007/gcas/p2p"
)

type MessageStoreFile struct {
	ID   string
	Key  string
	Size int64
}

type MessageGetFile struct {
	ID  string
	Key string
}

// GetFile is used to get the file from the local storage (file server) and broadcast it to the peers is not found locally
func (fs *FileServer) GetFile(key string) (io.Reader, error) {
	// check if the file with key already exists in the local storage
	if fs.Storage.Has(fs.Config.ID, key) {
		fmt.Printf("[%s] Serving file %s from local storage\n", fs.Config.Transport.Address(), key)
		_, r, err := fs.Storage.Read(fs.Config.ID, key)
		return r, err
	}
	fmt.Printf("[%s] File %s not found in local storage, fetching it from the network now...\n", fs.Config.Transport.Address(), key)
	msg := BroadCastMessage{
		Payload: MessageGetFile{
			ID:  fs.Config.ID,
			Key: cascrypto.HashKey(key),
		},
	}
	// broadcast the message to the peers
	if err := fs.BroadCastMessageToPeers(&msg); err != nil {
		return nil, err
	}
	time.Sleep(time.Millisecond * 500)
	// receive the file from the peers
	for _, peer := range fs.peers {
		// read the file size so we can limit the byte size
		var fileSize int64
		binary.Read(peer, binary.LittleEndian, &fileSize)

		n, err := fs.Storage.WriteDecrypt(fs.Config.EncryptionKey, fs.Config.ID, key, io.LimitReader(peer, fileSize))
		if err != nil {
			return nil, err
		}

		fmt.Printf("[%s] received (%d) bytes over the network from (%s)", fs.Config.Transport.Address(), n, peer.Address())

		peer.CloseStream()
	}
	// read the file from the local storage
	_, r, err := fs.Storage.Read(fs.Config.ID, key)
	return r, err
}

// StoreFile is used to store the file in the local storage (file server) and broadcast it to the peers
func (fs *FileServer) StoreFile(key string, r io.Reader) error {
	var (
		fileBuffer = new(bytes.Buffer)
		tee        = io.TeeReader(r, fileBuffer)
	)
	size, err := fs.Storage.Write(fs.Config.ID, key, tee)
	if err != nil {
		return err
	}
	msg := BroadCastMessage{
		Payload: MessageStoreFile{
			ID:   fs.Config.ID,
			Key:  cascrypto.HashKey(key),
			Size: size + 16,
		},
	}
	// broadcast the message to the peers
	if err := fs.BroadCastMessageToPeers(&msg); err != nil {
		return err
	}
	// delay to allow the peers to receive the message
	time.Sleep(time.Millisecond * 500)
	peers := []io.Writer{}
	for _, peer := range fs.peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	mw.Write([]byte{p2p.IncomingMessage})
	n, err := cascrypto.CopyEncrypt(fs.Config.EncryptionKey, fileBuffer, mw)
	if err != nil {
		return err
	}

	fmt.Printf("[%s] received and written (%d) bytes to disk\n", fs.Config.Transport.Address(), n)

	return nil
}
