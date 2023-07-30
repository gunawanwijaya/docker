package z

import (
	"crypto/rand"
	"crypto/subtle"
	"io"

	"golang.org/x/crypto/nacl/auth"
	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
	"golang.org/x/crypto/nacl/sign"
)

type keygen_nacl struct{}

func (keygen_nacl) Sign(r io.Reader) (privateKey *[2 * NACL_KEY_SIZE]byte, err error) {
	_, key, err := sign.GenerateKey(r)
	return key, err
}

func (keygen_nacl) Box(r io.Reader) (publicKey, privateKey *[NACL_KEY_SIZE]byte, err error) {
	return box.GenerateKey(r)
}

const (
	NACL_KEY_SIZE   int = 32
	NACL_NONCE_SIZE int = 24
)

type crypting_nacl struct{}

// ---------------------------------------------------------------------------------------

func (crypting_nacl) SecretKey(secretKey *[NACL_KEY_SIZE]byte) Crypter {
	return crypting_nacl_secret_key{secretKey}
}

type crypting_nacl_secret_key struct {
	secretKey *[NACL_KEY_SIZE]byte
}

func (x crypting_nacl_secret_key) Encrypt(msg []byte) (c CipherText, err error) {
	nonce := util{}.nacl.Nonce()
	p := secretbox.Seal(nonce[:], msg, &nonce, x.secretKey)
	return cipherText{Text(p)}, err
}

func (x crypting_nacl_secret_key) Decrypt(c CipherText) (msg []byte, err error) {
	nonce, p := util{}.nacl.NonceAndBodyOf(c.Bytes())
	msg, ok := secretbox.Open(nil, p, &nonce, x.secretKey)
	return util{}.nacl.OpenResultToErr(msg, ok)
}

// ---------------------------------------------------------------------------------------

func (crypting_nacl) PublicKey(peersPublicKey, privateKey *[NACL_KEY_SIZE]byte) Crypter {
	return crypting_nacl_public_key{peersPublicKey, privateKey}
}

type crypting_nacl_public_key struct {
	peersPublicKey *[NACL_KEY_SIZE]byte
	privateKey     *[NACL_KEY_SIZE]byte
}

func (x crypting_nacl_public_key) Encrypt(msg []byte) (c CipherText, err error) {
	nonce := util{}.nacl.Nonce()
	p := box.Seal(nonce[:], msg, &nonce, x.peersPublicKey, x.privateKey)
	return cipherText{Text(p)}, err
}

func (x crypting_nacl_public_key) Decrypt(c CipherText) (msg []byte, err error) {
	nonce, p := util{}.nacl.NonceAndBodyOf(c.Bytes())
	msg, ok := box.Open(nil, p, &nonce, x.peersPublicKey, x.privateKey)
	return util{}.nacl.OpenResultToErr(msg, ok)
}

// ---------------------------------------------------------------------------------------

func (crypting_nacl) SharedKey(sharedKey *[NACL_KEY_SIZE]byte) Crypter {
	return crypting_nacl_shared_key{sharedKey}
}

type crypting_nacl_shared_key struct {
	sharedKey *[NACL_KEY_SIZE]byte
}

func (x crypting_nacl_shared_key) Encrypt(msg []byte) (c CipherText, err error) {
	nonce := util{}.nacl.Nonce()
	p := box.SealAfterPrecomputation(nonce[:], msg, &nonce, x.sharedKey)
	return cipherText{Text(p)}, err
}

func (x crypting_nacl_shared_key) Decrypt(c CipherText) (msg []byte, err error) {
	nonce, p := util{}.nacl.NonceAndBodyOf(c.Bytes())
	msg, ok := box.OpenAfterPrecomputation(nil, p, &nonce, x.sharedKey)
	return util{}.nacl.OpenResultToErr(msg, ok)
}

// ---------------------------------------------------------------------------------------

func (crypting_nacl) Anonymous(peersPublicKey, privateKey, publicKey *[NACL_KEY_SIZE]byte) Crypter {
	return crypting_nacl_anonymous{peersPublicKey, privateKey, publicKey}
}

type crypting_nacl_anonymous struct {
	peersPublicKey *[NACL_KEY_SIZE]byte
	privateKey     *[NACL_KEY_SIZE]byte
	publicKey      *[NACL_KEY_SIZE]byte
}

func (x crypting_nacl_anonymous) Encrypt(msg []byte) (c CipherText, err error) {
	p, err := box.SealAnonymous(nil, msg, x.peersPublicKey, rand.Reader)
	return cipherText{Text(p)}, err
}

func (x crypting_nacl_anonymous) Decrypt(c CipherText) (msg []byte, err error) {
	msg, ok := box.OpenAnonymous(nil, c.Bytes(), x.publicKey, x.privateKey)
	return util{}.nacl.OpenResultToErr(msg, ok)
}

// ---------------------------------------------------------------------------------------

func (signing) NaCl(privateKey *[2 * NACL_KEY_SIZE]byte) Signer {
	publicKey := new([NACL_KEY_SIZE]byte)
	_ = copy(publicKey[:], privateKey[NACL_KEY_SIZE:])
	return signing_nacl{privateKey, publicKey}
}

type signing_nacl struct {
	privateKey *[2 * NACL_KEY_SIZE]byte
	publicKey  *[NACL_KEY_SIZE]byte
}

func (x signing_nacl) Sign(msg []byte) (s Signature, err error) {
	p := sign.Sign(nil, msg, x.privateKey)
	return signature{Text(p)}, nil
}

func (x signing_nacl) Verify(s Signature, msg []byte) (err error) {
	if p, ok := sign.Open(nil, msg, x.publicKey); !ok || subtle.ConstantTimeCompare(p, msg) != 1 {
		err = ErrInvalidSignature
	}
	return err
}

// ---------------------------------------------------------------------------------------

func (hashing) NaCl(secretKey *[NACL_KEY_SIZE]byte) Hasher {
	return hashing_nacl{secretKey}
}

type hashing_nacl struct {
	secretKey *[NACL_KEY_SIZE]byte
}

func (x hashing_nacl) Hash(key []byte) (enc HashedKey, err error) {
	p := auth.Sum(key, x.secretKey)[:]
	return hashedKey{Text(p)}, nil
}

func (x hashing_nacl) Compare(enc HashedKey, key []byte) error {
	if !auth.Verify(enc.Bytes(), key, x.secretKey) {
		return ErrKeyComparisonFailed
	}
	return nil
}
