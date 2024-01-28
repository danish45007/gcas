package cas

import (
	"bytes"
	"encoding/gob"

	"github.com/danish45007/gcas/p2p"
)

type BroadCastMessage struct {
	Payload any
}

func (fs *FileServer) BroadCastMessageToPeers(message *BroadCastMessage) error {
	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(message)
	if err != nil {
		return err
	}
	// loop over all the peer map that file server is maintaining
	for _, peer := range fs.peers {
		// send the incoming message byte
		peer.Send([]byte{p2p.IncomingMessage})
		// send the message payload
		err := peer.Send(buffer.Bytes())
		if err != nil {
			return err
		}
	}
	return nil
}
