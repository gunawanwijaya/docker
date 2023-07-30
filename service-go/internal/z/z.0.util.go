package z

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
)

// ---------------------------------------------------------------------------------------

type util struct {
	aes  aesUtil
	nacl naclUtil
}

func (keygen) Nonce(n int) []byte {
	p := make([]byte, n, n)
	_, _ = io.ReadFull(rand.Reader, p)
	return p
}

type naclUtil struct{}

func (naclUtil) OpenResultToErr(msg []byte, ok bool) ([]byte, error) {
	if !ok || msg == nil {
		return nil, ErrInvalidEncryption
	}
	return msg, nil
}

func (naclUtil) NewKey() *[32]byte {
	key := new([32]byte)
	_, _ = io.ReadFull(rand.Reader, key[:])
	return key
}

func (naclUtil) Nonce() (nonce [24]byte) {
	_, _ = io.ReadFull(rand.Reader, nonce[:])
	return nonce
}

func (naclUtil) NonceAndBodyOf(ciphertext []byte) (nonce [24]byte, p []byte) {
	_, p = copy(nonce[:], ciphertext[:24]), ciphertext[24:]
	return nonce, p
}

type aesUtil struct{}

func (aesUtil) ErrInvalidPadding() error   { return errors.New("invalid padding") }
func (aesUtil) ErrInvalidBlockSize() error { return errors.New("invalid blocksize") }
func (aesUtil) ErrInvalidPKCS7Data() error { return errors.New("invalid PKCS7 data") }

// PKCS7Padding right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func (x aesUtil) PKCS7Padding(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, x.ErrInvalidBlockSize()
	}
	if b == nil || len(b) == 0 {
		return nil, x.ErrInvalidPKCS7Data()
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// PKCS7Trimming validates and trims data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func (x aesUtil) PKCS7Trimming(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, x.ErrInvalidBlockSize()
	}
	if b == nil || len(b) == 0 {
		return nil, x.ErrInvalidPKCS7Data()
	}
	if len(b)%blocksize != 0 {
		return nil, x.ErrInvalidPadding()
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, x.ErrInvalidPadding()
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, x.ErrInvalidPadding()
		}
	}
	return b[:len(b)-n], nil
}

func (aesUtil) PKCS5Padding(b []byte, blockSize int) []byte {
	n := blockSize - len(b)%blockSize
	t := bytes.Repeat([]byte{byte(n)}, n)
	return append(b, t...)
}

func (aesUtil) PKCS5Trimming(b []byte) []byte {
	n := b[len(b)-1]
	return b[:len(b)-int(n)]
}

func (keygen) CryptoRandomReader() io.Reader { return rand.Reader }
