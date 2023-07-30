package z

import "errors"

var _ Texter = Text(nil)

type Text []byte

func (val Text) Bytes() []byte    { return val }
func (val Text) String() string   { return string(val) }
func (val Text) GoString() string { return string(val) }

type Texter interface {
	Bytes() []byte
	String() string   // fmt.Stringer
	GoString() string // fmt.GoStringer
}

var Strategy strategy

type strategy struct {
	Crypting crypting
	Encoding encoding
	Hashing  hashing
	Signing  signing

	Keygen keygen
	Build  build
}

type build struct{}

type keygen struct {
	NaCl keygen_nacl
}

type crypting struct {
	AES  crypting_aes
	NaCl crypting_nacl
	RSA  crypting_rsa
}

type encoding struct{}

type hashing struct {
	Password hashing_pw
}

type signing struct {
	RSA signing_rsa
}

type hashing_pw struct{}

var ErrInvalidEncryption = errors.New("invalid encryption")
var ErrInvalidSignature = errors.New("invalid signature")
var ErrKeyComparisonFailed = errors.New("key comparison failed")
var ErrInvalidPassword = errors.New("invalid password")

// ---------------------------------------------------------------------------------------
type Hasher interface {
	Hash(key []byte) (enc HashedKey, err error)
	Compare(enc HashedKey, key []byte) (err error)
}
type HashedKey interface {
	Texter
	Compare(s Hasher, key []byte) error
}
type hashedKey struct{ Texter }

func (x hashedKey) Compare(s Hasher, key []byte) error { return s.Compare(x, key) }
func (build) HashedKey(x Text) HashedKey               { return hashedKey{x} }

// ---------------------------------------------------------------------------------------
type Crypter interface {
	Encrypt(msg []byte) (c CipherText, err error)
	Decrypt(c CipherText) (msg []byte, err error)
}
type CipherText interface {
	Texter
	Decrypt(s Crypter) (msg []byte, err error)
}
type cipherText struct{ Texter }

func (x cipherText) Decrypt(s Crypter) (msg []byte, err error) { return s.Decrypt(x) }
func (build) CipherText(x Text) CipherText                     { return cipherText{x} }

// ---------------------------------------------------------------------------------------
type Signer interface {
	Sign(msg []byte) (s Signature, err error)
	Verify(s Signature, msg []byte) (err error)
}
type Signature interface {
	Texter
	Verify(s Signer, msg []byte) (err error)
}
type signature struct{ Texter }

func (x signature) Verify(s Signer, msg []byte) (err error) { return s.Verify(x, msg) }
func (build) Signature(x Text) Signature                    { return signature{x} }

// ---------------------------------------------------------------------------------------
type EncoderDecoder interface {
	Encode(msg []byte) (e EncodedValue, err error)
	Decode(e EncodedValue) (msg []byte, err error)
}
type EncodedValue interface {
	Texter
	Decode(e EncoderDecoder) (msg []byte, err error)
}
type encodedValue struct{ Texter }

func (x encodedValue) Decode(s EncoderDecoder) (msg []byte, err error) { return s.Decode(x) }
func (build) EncodedValue(x Text) EncodedValue                         { return encodedValue{x} }
