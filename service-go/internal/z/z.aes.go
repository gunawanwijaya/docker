package z

import (
	"crypto/aes"
	"crypto/cipher"
)

type crypting_aes struct{}

// ---------------------------------------------------------------------------------------

func (crypting_aes) CBC(key []byte, pkcsPaddingMode uint) Crypter {
	block, err := aes.NewCipher(key)
	return crypting_aes_cbc{pkcsPaddingMode, block, err}
}

type crypting_aes_cbc struct {
	pkcsPaddingMode uint

	block cipher.Block
	err   error
}

func (x crypting_aes_cbc) Encrypt(msg []byte) (c CipherText, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	switch x.pkcsPaddingMode {
	default:
		return nil, util{}.aes.ErrInvalidPadding()
	case 5:
		msg = util{}.aes.PKCS5Padding(msg, aes.BlockSize)
	case 7:
		msg, err = util{}.aes.PKCS7Padding(msg, aes.BlockSize)
		if err != nil {
			return nil, err
		}
	}
	iv := keygen{}.Nonce(x.block.BlockSize())
	enc := make([]byte, len(msg))
	cipher.NewCBCEncrypter(x.block, iv).CryptBlocks(enc, msg)
	return cipherText{Text(append(iv, enc...))}, nil
}

func (x crypting_aes_cbc) Decrypt(c CipherText) (msg []byte, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	cc, n := c.Bytes(), x.block.BlockSize()
	iv, enc := cc[:n], cc[n:]
	msg = make([]byte, len(enc))
	cipher.NewCBCDecrypter(x.block, iv).CryptBlocks(msg, enc)
	switch x.pkcsPaddingMode {
	default:
		return nil, util{}.aes.ErrInvalidPadding()
	case 5:
		msg = util{}.aes.PKCS5Trimming(msg)
	case 7:
		msg, err = util{}.aes.PKCS7Trimming(msg, aes.BlockSize)
		if err != nil {
			return nil, err
		}
	}
	return msg, nil
}

// ---------------------------------------------------------------------------------------

func (crypting_aes) CFB(key []byte) Crypter {
	block, err := aes.NewCipher(key)
	return crypting_aes_cfb{block, err}
}

type crypting_aes_cfb struct {
	block cipher.Block
	err   error
}

func (x crypting_aes_cfb) Encrypt(msg []byte) (c CipherText, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	iv := keygen{}.Nonce(x.block.BlockSize())
	enc := make([]byte, len(msg))
	cipher.NewCFBEncrypter(x.block, iv).XORKeyStream(enc, msg)
	return cipherText{Text(append(iv, enc...))}, nil
}

func (x crypting_aes_cfb) Decrypt(c CipherText) (msg []byte, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	cc, n := c.Bytes(), x.block.BlockSize()
	iv, enc := cc[:n], cc[n:]
	msg = make([]byte, len(enc))
	cipher.NewCFBDecrypter(x.block, iv).XORKeyStream(msg, enc)
	return msg, nil
}

// ---------------------------------------------------------------------------------------

func (crypting_aes) CTR(key []byte) Crypter {
	block, err := aes.NewCipher(key)
	return crypting_aes_ctr{block, err}
}

type crypting_aes_ctr struct {
	block cipher.Block
	err   error
}

func (x crypting_aes_ctr) Encrypt(msg []byte) (c CipherText, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	iv, enc := keygen{}.Nonce(x.block.BlockSize()), make([]byte, len(msg))
	cipher.NewCTR(x.block, iv).XORKeyStream(enc, msg)
	return cipherText{Text(append(iv, enc...))}, nil
}

func (x crypting_aes_ctr) Decrypt(c CipherText) (msg []byte, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	cc, n := c.Bytes(), x.block.BlockSize()
	iv, enc := cc[:n], cc[n:]
	msg = make([]byte, len(enc))
	cipher.NewCTR(x.block, iv).XORKeyStream(msg, enc)
	return msg, nil
}

// ---------------------------------------------------------------------------------------

func (crypting_aes) GCM(key []byte, additionalData []byte) Crypter {
	block, err := aes.NewCipher(key)
	if err != nil {
		return crypting_aes_gcm{nil, nil, err}
	}
	gcm, err := cipher.NewGCM(block)
	return crypting_aes_gcm{additionalData, gcm, err}
}

type crypting_aes_gcm struct {
	additionalData []byte

	gcm cipher.AEAD
	err error
}

func (x crypting_aes_gcm) Encrypt(msg []byte) (c CipherText, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	iv := keygen{}.Nonce(x.gcm.NonceSize())
	enc := x.gcm.Seal(iv, iv, msg, x.additionalData)
	return cipherText{Text(enc)}, nil
}

func (x crypting_aes_gcm) Decrypt(c CipherText) (msg []byte, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	cc, n := c.Bytes(), x.gcm.NonceSize()
	iv, enc := cc[:n], cc[n:]
	msg, err = x.gcm.Open(nil, iv, enc, x.additionalData)
	return msg, err
}

// ---------------------------------------------------------------------------------------

func (crypting_aes) OFB(key []byte) Crypter {
	block, err := aes.NewCipher(key)
	return crypting_aes_ofb{block, err}
}

type crypting_aes_ofb struct {
	block cipher.Block
	err   error
}

func (x crypting_aes_ofb) Encrypt(msg []byte) (c CipherText, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	iv, enc := keygen{}.Nonce(x.block.BlockSize()), make([]byte, len(msg))
	cipher.NewOFB(x.block, iv).XORKeyStream(enc, msg)
	return cipherText{Text(append(iv, enc...))}, nil
}

func (x crypting_aes_ofb) Decrypt(c CipherText) (msg []byte, err error) {
	if err = x.err; err != nil {
		return nil, err
	}
	cc, n := c.Bytes(), x.block.BlockSize()
	iv, enc := cc[:n], cc[n:]
	msg = make([]byte, len(enc))
	cipher.NewOFB(x.block, iv).XORKeyStream(msg, enc)
	return msg, nil
}
