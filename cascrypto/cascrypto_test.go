package cascrypto

import (
	"bytes"
	"testing"
)

func TestCopyEncryptionDecryption(t *testing.T) {
	payload := []byte("test data")
	// source reader with payload
	source := bytes.NewReader([]byte(payload))
	// destination writer buffer
	destination := new(bytes.Buffer)
	// runtime encryption key
	key := NewEncryptionKey()
	_, err := CopyEncrypt(key, source, destination)
	if err != nil {
		t.Error(err)
	}

	out := new(bytes.Buffer)
	bytesWritten, err := CopyDecrypt(key, destination, out)
	if err != nil {
		t.Error(err)
	}

	if bytesWritten != 16+len(payload) {
		t.Fail()
	}

	if !bytes.Equal(out.Bytes(), payload) {
		t.Fail()
	}

}
