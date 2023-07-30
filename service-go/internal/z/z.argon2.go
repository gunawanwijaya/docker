package z

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// ---------------------------------------------------------------------------------------

// Argon2 default arguments:
//
//	mode = 2       // 0 = argon2d | 1 = argon2i | 2 = argon2id
//	t    = 1       // (t)imes or number of iterations
//	m    = 64*1024 // (m)emory in kilobytes
//	p    = 4       // (p)arallelism or number of threads
//	salt = []byte  // random
//	n    = 32      // le(n)gth of hash result
func (hashing_pw) Argon2(mode uint8, t uint32, m uint32, p uint8, salt []byte, n int) Hasher {
	return hashing_pw_argon2{mode, t, m, p, salt, n}
}

type hashing_pw_argon2 struct {
	mode uint8  // mode : 0 = d (unsupported), 1 = i, 2 = id
	t    uint32 // t    : times (iteration)
	m    uint32 // m    : memory in KB
	p    uint8  // p    : parallelism (number of threads)
	salt []byte // salt
	n    int    // n    : length of hash result
}

func (x hashing_pw_argon2) Hash(key []byte) (HashedKey, error) {
	kdf := argon2.IDKey
	if x.mode == 1 {
		kdf = argon2.Key
	}
	modeStr := map[uint8]string{
		0: "argon2d",
		1: "argon2i",
		2: "argon2id",
	}[x.mode]
	salt, _ := b64.Encode(x.salt)
	rest, _ := b64.Encode(kdf(key, x.salt, x.t, x.m, x.p, uint32(x.n)))
	p := fmt.Sprintf("$%s$v=%d$t=%d,m=%d,p=%d$%s$%s",
		modeStr, argon2.Version, x.t, x.m, x.p,
		salt.Bytes(),
		rest,
	)

	return hashedKey{Text(p)}, nil
}

func (hashing_pw_argon2) Compare(enc HashedKey, key []byte) error {
	x, err := new(hashing_pw_argon2).decode(enc.String())
	if err != nil {
		return err
	}
	h, err := x.Hash(key)
	if err != nil {
		return err
	}
	if subtle.ConstantTimeCompare(enc.Bytes(), h.Bytes()) != 1 {
		return ErrKeyComparisonFailed
	}
	return nil
}

func (x hashing_pw_argon2) decode(text string) (_ hashing_pw_argon2, err error) {
	var next int
	err = errors.New("argon2: invalid encoded format")
	switch strings.ToLower(text[:9]) {
	default:
		return x, err
	case "$argon2d$":
		next, x.mode = 8, 0
	case "$argon2i$":
		next, x.mode = 8, 1
	case "$argon2id":
		next, x.mode = 9, 2
	}

	if text[next:next+6] != "$v=19$" {
		return x, err
	}
	next = next + 6

	for stop := false; !stop && next < len(text); {
		var err0 error
		if next, stop, err0 = x.extractTMP(text, &x.t, &x.m, &x.p, next); err0 != nil {
			return x, err
		}
	}

	d := strings.IndexRune(text[next:], '$')
	if d < 0 {
		return x, err
	}
	salt, _ := encodedValue{Text(text[next : next+d])}.Decode(b64)
	rest, _ := encodedValue{Text(text[next+d+1:])}.Decode(b64)
	x.salt = salt
	x.n = len(rest)
	valid := (x.mode == 1 || x.mode == 2) &&
		x.t > 0 &&
		x.m > 0 &&
		x.p > 0 &&
		len(x.salt) > 0 &&
		x.n > 0
	if !valid {
		return x, err
	}
	return x, nil
}

func (hashing_pw_argon2) extractTMP(text string, t, m *uint32, p *uint8, next int) (_ int, stop bool, err error) {
	digitEnd := len(text) - 1
	c := strings.IndexRune(text[next+2:], ',')
	d := strings.IndexRune(text[next+2:], '$')
	if c > 0 {
		digitEnd, stop = c+next+2, false
	} else if d > 0 && c < 0 {
		digitEnd, stop = d+next+2, true
	} else if c < 0 && d < 0 {
		stop = true
	}

	digit, dstr := 0, string(text[next+2:digitEnd])
	digit, err = strconv.Atoi(dstr)

	switch strings.ToLower(text[next : next+2]) {
	case "m=":
		*m = uint32(digit)
	case "p=":
		*p = uint8(digit)
	case "t=":
		*t = uint32(digit)
	}
	next = digitEnd + 1
	return next, stop, err
}

// ---------------------------------------------------------------------------------------

var (
	b64 = Strategy.Encoding.Base64(0)
)
