package storage

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

type PathTransformFunc func(string) Path

func DefaultPathTransform(key string) Path {
	return Path{
		PathName: key,
		FileName: key,
	}
}
func CASPathTransform(key string) Path {
	blockSize := 5
	hash := sha1.Sum([]byte(key))
	hashString := hex.EncodeToString(hash[:])
	sliceLength := len(hashString) / blockSize
	paths := make([]string, sliceLength)
	for i := 0; i < sliceLength; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashString[from:to]
	}
	return Path{
		PathName: strings.Join(paths, "/"),
		FileName: hashString,
	}
}
