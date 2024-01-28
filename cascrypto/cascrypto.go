package cascrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

func GenerateId() string {
	buffer := make([]byte, 32)
	io.ReadFull(rand.Reader, buffer)
	return hex.EncodeToString(buffer)
}

func HashKey(key string) string {
	md5Hash := md5.Sum([]byte(key))
	return hex.EncodeToString(md5Hash[:])
}

func NewEncryptionKey() []byte {
	keyBuffer := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuffer)
	return keyBuffer
}

func copyStream(stream cipher.Stream, blockSize int, dst io.Writer, src io.Reader) (int, error) {
	var (
		// Buffer to hold chunks of data read from the source
		buffer = make([]byte, 32*1024)
		// Initialize with the block size (e.g., IV size for encryption)
		totalBytesWritten = blockSize
	)
	// Loop until the end of the source data is reached
	for {
		// Read a chunk of data from the source into the buffer
		n, err := src.Read(buffer)
		if n > 0 {
			// XOR the chunk with the keystream obtained from the cipher.Stream
			stream.XORKeyStream(buffer, buffer[:n])

			// Write the XORed data to the destination
			bytesWritten, err := dst.Write(buffer[:n])
			if err != nil {
				return 0, err
			}

			// Update the total number of bytes written
			totalBytesWritten += bytesWritten
		}

		// Check for the end of the source data
		if err == io.EOF {
			break
		}

		// Propagate any other errors that occurred during reading
		if err != nil {
			return 0, err
		}
	}

	return totalBytesWritten, nil
}

func CopyDecrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	// Read the IV from the given io.Reader which, in our case should be the
	// the block.BlockSize() bytes we read.
	iv := make([]byte, block.BlockSize())
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), dst, src)
}

func CopyEncrypt(key []byte, src io.Reader, dst io.Writer) (int, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize()) // 16 bytes
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepend the IV to the file.
	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, block.BlockSize(), dst, src)
}
