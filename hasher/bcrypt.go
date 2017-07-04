package hasher

import "golang.org/x/crypto/bcrypt"

// bcryptHasher is a Hasher that uses bcrypt under the hood.
type bcryptHasher struct {
	Cost int
}

// NewBcryptHasher returns a new bcryptHasher.
func NewBcryptHasher(cost int) Hasher {
	return bcryptHasher{cost}
}

func (h bcryptHasher) GenerateHashFromPassword(password []byte) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(password), h.Cost)
}

func (h bcryptHasher) CompareHashWithPassword(hash, password []byte) (err error) {
	return bcrypt.CompareHashAndPassword(hash, password)
}
