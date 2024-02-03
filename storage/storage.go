package storage

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/danish45007/gcas/cascrypto"
)

const defaultRootFolderName = "gcasnetworks"

type StorageConfig struct {
	PathTransform PathTransformFunc
	SetRoot       string
}

type Storage struct {
	config StorageConfig
}

func NewStorage(config StorageConfig) *Storage {
	if config.PathTransform == nil {
		config.PathTransform = DefaultPathTransform
	}
	if len(config.SetRoot) == 0 {
		config.SetRoot = defaultRootFolderName
	}
	return &Storage{
		config: config,
	}
}

func (s *Storage) Has(id, key string) bool {
	pathKey := s.config.PathTransform(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.config.SetRoot, id, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Storage) Clear() error {
	return os.RemoveAll(s.config.SetRoot)
}

func (s *Storage) Delete(id, key string) error {
	pathKey := s.config.PathTransform(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.config.SetRoot, id, pathKey.FullPath())
	defer func() {
		log.Printf("Deleted %s from disk", fullPathWithRoot)
	}()

	err := os.RemoveAll(fullPathWithRoot)
	if err != nil {
		return err
	}

	// After deletion, traverse upward and delete empty directories
	err = s.deleteEmptyDirectories(filepath.Dir(fullPathWithRoot))
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) deleteEmptyDirectories(dirPath string) error {
	for dirPath != s.config.SetRoot && isEmptyDirectory(dirPath) {
		err := os.Remove(dirPath)
		if err != nil {
			return err
		}
		dirPath = filepath.Dir(dirPath)
	}
	return nil
}

func isEmptyDirectory(dirPath string) bool {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		log.Printf("Error reading directory %s: %s", dirPath, err)
		return false
	}

	return len(dir) == 0
}

func (s *Storage) Read(id, key string) (int64, io.Reader, error) {
	return s.ReadStream(id, key)
}

func (s *Storage) ReadStream(id, key string) (int64, io.ReadCloser, error) {
	pathKey := s.config.PathTransform(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.config.SetRoot, id, pathKey.FullPath())
	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}
	return fileInfo.Size(), file, nil
}

func (s *Storage) openFileForWriting(id, key string) (*os.File, error) {
	pathKey := s.config.PathTransform(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.config.SetRoot, id, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil, err
	}
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.config.SetRoot, id, pathKey.FullPath())
	return os.Create(fullPathWithRoot)
}

func (s *Storage) WriteDecrypt(encKey []byte, id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	n, err := cascrypto.CopyDecrypt(encKey, r, f)
	return int64(n), err
}

func (s *Storage) WriteStream(id string, key string, r io.Reader) (int64, error) {
	file, err := s.openFileForWriting(id, key)
	if err != nil {
		return 0, err
	}
	return io.Copy(file, r)

}

func (s *Storage) Write(id string, key string, r io.Reader) (int64, error) {
	return s.WriteStream(id, key, r)
}
