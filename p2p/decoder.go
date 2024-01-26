package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(r io.Reader, message *Message) error
}

type DefaultDecoder struct{}
type GoBinaryDecoder struct{}

func (d DefaultDecoder) Decode(r io.Reader, message *Message) error {
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	message.Payload = buf[:n]
	return nil
}

func (d GoBinaryDecoder) Decode(r io.Reader, message *Message) error {
	return gob.NewDecoder(r).Decode(message)
}
