package storage

import (
	"bytes"
	"testing"
)

// func TestStorage(t *testing.T) {
// 	config := StorageConfig{
// 		PathTransform: CASPathTransform,
// 		SetRoot:       "../",
// 	}
// 	testKey := "test"
// 	testData := []byte("test data...")
// 	s := NewStorage(config)
// 	data := bytes.NewReader(testData)
// 	if err := s.WriteStream(testKey, data); err != nil {
// 		t.Error(err)
// 	}
// 	r, err := s.Read(testKey)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	b, _ := io.ReadAll(r)
// 	if !bytes.Equal(b, testData) {
// 		t.Errorf("Expected %s, got %s", testData, b)
// 	}
// }

func TestStorageDelete(t *testing.T) {
	config := StorageConfig{
		PathTransform: CASPathTransform,
		SetRoot:       "../",
	}
	testKey := "test"
	testData := []byte("test data...")
	s := NewStorage(config)
	data := bytes.NewReader(testData)
	if err := s.WriteStream(testKey, data); err != nil {
		t.Error(err)
	}
	if err := s.Delete(testKey); err != nil {
		t.Error(err)
	}
	if s.Has(testKey) {
		t.Errorf("Expected key %s to be deleted", testKey)
	}
}
