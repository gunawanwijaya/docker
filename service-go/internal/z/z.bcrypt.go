package z

import "golang.org/x/crypto/bcrypt"

func (hashing_pw) Bcrypt(cost int) Hasher {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		cost = 10
	}
	return hashing_pw_bcrypt{cost, nil}
}

type hashing_pw_bcrypt struct {
	cost int    // cost
	salt []byte // salt
}

func (x hashing_pw_bcrypt) Hash(key []byte) (HashedKey, error) {
	p, err := bcrypt.GenerateFromPassword(key, x.cost)
	if err != nil {
		return nil, err
	}
	return hashedKey{Text(p)}, nil
}

func (hashing_pw_bcrypt) Compare(enc HashedKey, key []byte) error {
	if bcrypt.CompareHashAndPassword(enc.Bytes(), key) != nil {
		return ErrKeyComparisonFailed
	}
	return nil
}
