package z

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"io"
)

func (keygen) ECDSA(r io.Reader, curve int) (*ecdsa.PrivateKey, error) {
	var c elliptic.Curve
	switch curve {
	case 224:
		c = elliptic.P224()
	default:
		fallthrough
	case 256:
		c = elliptic.P256()
	case 384:
		c = elliptic.P384()
	case 521:
		c = elliptic.P521()
	}
	return ecdsa.GenerateKey(c, r)
}

// ---------------------------------------------------------------------------------------

func (signing) ECDSA(key *ecdsa.PrivateKey) Signer {
	return signECDSA{key}
}

type signECDSA struct {
	key *ecdsa.PrivateKey
}

func (x signECDSA) Sign(msg []byte) (s Signature, err error) {
	p, err := ecdsa.SignASN1(rand.Reader, x.key, msg)
	return signature{Text(p)}, err
}

func (x signECDSA) Verify(s Signature, msg []byte) (err error) {
	if ok := ecdsa.VerifyASN1(&x.key.PublicKey, msg, s.Bytes()); !ok {
		err = ErrInvalidSignature
	}
	return err
}
