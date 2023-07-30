package z

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

func (hashing) SHA(mode int, key []byte) Hasher {
	var h func() hash.Hash
	switch mode {
	default:
		h = sha256.New
	case 1:
		h = sha1.New
	case 224:
		h = sha256.New224
	case 256:
		h = sha256.New
	case 384:
		h = sha512.New384
	case 512:
		h = sha512.New
	case 512_224:
		h = sha512.New512_224
	case 512_256:
		h = sha512.New512_256
	}
	return hashing_sha{h, key}
}

type hashing_sha struct {
	h   func() hash.Hash
	key []byte
}

func (x hashing_sha) Hash(m []byte) (enc HashedKey, err error) {
	mac := hmac.New(x.h, x.key)
	_, _ = mac.Write(m)
	out := mac.Sum(nil)
	return hashedKey{Text(out)}, nil
}

func (x hashing_sha) Compare(enc HashedKey, key []byte) error {
	if len(key) != 32 {
		return ErrKeyComparisonFailed
	}
	mac := hmac.New(x.h, x.key)
	_, _ = mac.Write(key)
	exp := mac.Sum(nil)
	if !hmac.Equal(enc.Bytes(), exp) {
		return ErrKeyComparisonFailed
	}
	return nil
}
