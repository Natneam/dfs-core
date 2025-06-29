package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

func NewEncryptionKey() []byte {
	keyBuf := make([]byte, 32)
	io.ReadFull(rand.Reader, keyBuf)
	return keyBuf
}

func copyStream(stream cipher.Stream, blockSize int64, src io.Reader, dst io.Writer) (int64, error) {
	var (
		buf       = make([]byte, 32*1024)
		totalSize = blockSize
	)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			nn, err := dst.Write(buf[:n])
			if err != nil {
				return 0, err
			}
			totalSize += int64(nn)
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}
	}

	return totalSize, nil
}

func CopyEncrypt(key []byte, src io.Reader, dst io.Writer) (int64, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return 0, err
	}

	// prepend the IV to the stream
	if _, err := dst.Write(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)

	return copyStream(stream, int64(block.BlockSize()), src, dst)

}

func CopyDecrypt(key []byte, src io.Reader, dst io.Writer) (int64, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return 0, err
	}

	iv := make([]byte, block.BlockSize())
	if _, err := src.Read(iv); err != nil {
		return 0, err
	}

	stream := cipher.NewCTR(block, iv)
	return copyStream(stream, int64(block.BlockSize()), src, dst)
}
