package z

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"hash"
	"io"
)

func (keygen) RSA(r io.Reader, bits int) (*rsa.PrivateKey, error) {
	switch bits {
	case 2048, 3072, 4096:
	default:
		bits = 2048
	}
	return rsa.GenerateKey(r, bits)
}

type crypting_rsa struct{}

// ---------------------------------------------------------------------------------------

func (crypting_rsa) PKCS1v15(key *rsa.PrivateKey) Crypter {
	return crypting_rsa_pkcs1v15{key}
}

type crypting_rsa_pkcs1v15 struct {
	key *rsa.PrivateKey
}

func (x crypting_rsa_pkcs1v15) Encrypt(msg []byte) (c CipherText, err error) {
	p, err := rsa.EncryptPKCS1v15(rand.Reader, &x.key.PublicKey, msg)
	return cipherText{Text(p)}, err
}
func (x crypting_rsa_pkcs1v15) Decrypt(c CipherText) (msg []byte, err error) {
	p, err := rsa.DecryptPKCS1v15(nil, x.key, c.Bytes())
	return p, err
}

// ---------------------------------------------------------------------------------------

func (crypting_rsa) OAEP(hash hash.Hash, key *rsa.PrivateKey, label []byte) Crypter {
	return crypting_rsa_oaep{hash, key, label}
}

type crypting_rsa_oaep struct {
	hash  hash.Hash
	key   *rsa.PrivateKey
	label []byte
}

func (x crypting_rsa_oaep) Encrypt(msg []byte) (c CipherText, err error) {
	p, err := rsa.EncryptOAEP(x.hash, rand.Reader, &x.key.PublicKey, msg, x.label)
	return cipherText{Text(p)}, err
}

func (x crypting_rsa_oaep) Decrypt(c CipherText) (msg []byte, err error) {
	p, err := rsa.DecryptOAEP(x.hash, nil, x.key, c.Bytes(), x.label)
	return p, err
}

// ---------------------------------------------------------------------------------------

type signing_rsa struct{}

// ---------------------------------------------------------------------------------------

func (signing_rsa) PSS(hash crypto.Hash, key *rsa.PrivateKey, opts *rsa.PSSOptions) Signer {
	return signing_rsa_pss{hash, key, opts}
}

type signing_rsa_pss struct {
	hash crypto.Hash
	key  *rsa.PrivateKey
	opts *rsa.PSSOptions
}

func (x signing_rsa_pss) Sign(msg []byte) (s Signature, err error) {
	p, err := rsa.SignPSS(rand.Reader, x.key, x.hash, msg, x.opts)
	return signature{Text(p)}, err
}

func (x signing_rsa_pss) Verify(s Signature, msg []byte) (err error) {
	err = rsa.VerifyPSS(&x.key.PublicKey, x.hash, msg, s.Bytes(), x.opts)
	return err
}

// ---------------------------------------------------------------------------------------

func (signing_rsa) PKCS1v15(hash crypto.Hash, key *rsa.PrivateKey) Signer {
	return signing_rsa_pkcs1v15{hash, key}
}

type signing_rsa_pkcs1v15 struct {
	hash crypto.Hash
	key  *rsa.PrivateKey
}

func (x signing_rsa_pkcs1v15) Sign(msg []byte) (s Signature, err error) {
	p, err := rsa.SignPKCS1v15(nil, x.key, x.hash, msg)
	return signature{Text(p)}, err
}

func (x signing_rsa_pkcs1v15) Verify(s Signature, msg []byte) (err error) {
	err = rsa.VerifyPKCS1v15(&x.key.PublicKey, x.hash, msg, s.Bytes())
	return err
}
