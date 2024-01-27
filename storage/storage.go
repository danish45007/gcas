package storage

import (
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
)

type StorageConfig struct {
	PathTransform PathTransformFunc
	SetRoot       string
}

type Storage struct {
	config StorageConfig
}

func NewStorage(config StorageConfig) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Has(key string) bool {
	pathKey := s.config.PathTransform(key)
	fileNameWithPath := s.config.SetRoot + pathKey.FullPath()
	_, err := os.Stat(fileNameWithPath)
	return err == nil
}

func (s *Storage) Delete(key string) error {
	pathKey := s.config.PathTransform(key)
	fileNameWithPath := s.config.SetRoot + pathKey.FullPath()
	defer func() {
		log.Printf("Deleted %s from disk", fileNameWithPath)
	}()

	err := os.RemoveAll(fileNameWithPath)
	if err != nil {
		return err
	}

	// After deletion, traverse upward and delete empty directories
	err = s.deleteEmptyDirectories(filepath.Dir(fileNameWithPath))
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

func (s *Storage) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	if err != nil {
		return nil, err
	}
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	defer f.Close()
	return buf, err
}

func (s *Storage) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.config.PathTransform(key)
	fileNameWithPath := s.config.SetRoot + pathKey.FullPath()
	return os.Open(fileNameWithPath)
}

func (s *Storage) WriteStream(key string, r io.Reader) error {
	pathKey := s.config.PathTransform(key)
	// create parent directory
	// create a dir a level out from current dir
	if err := os.MkdirAll(s.config.SetRoot+pathKey.PathName, os.ModePerm); err != nil {
		return err
	}
	// create file
	fileNameWithPath := s.config.SetRoot + pathKey.FullPath()

	file, err := os.Create(fileNameWithPath)
	if err != nil {
		return err
	}
	bytes, err := io.Copy(file, r)
	if err != nil {
		return err
	}
	log.Printf("Wrote %d bytes to disk %s", bytes, fileNameWithPath)
	return nil

}
