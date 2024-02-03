package storage

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/danish45007/gcas/cascrypto"
)

func newStore() *Storage {
	opts := StorageConfig{
		PathTransform: CASPathTransform,
	}
	return NewStorage(opts)
}

func teardown(t *testing.T, s *Storage) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
func TestStorage(t *testing.T) {
	newStorage := newStore()
	id := cascrypto.GenerateId()
	defer teardown(t, newStorage)
	i := 1
	key := fmt.Sprintf("test_%d", i)
	data := []byte(fmt.Sprintf("test data %d", i))
	_, err := newStorage.WriteStream(id, key, bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}
	if exist := newStorage.Has(id, key); !exist {
		t.Errorf("expected to have key %s", key)
	}
	_, r, err := newStorage.Read(id, key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)
	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}

	if err := newStorage.Delete(id, key); err != nil {
		t.Error(err)
	}

	if ok := newStorage.Has(id, key); ok {
		t.Errorf("expected to NOT have key %s", key)
	}
}
