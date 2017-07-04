package hasher

// Hasher can generate crypto hashes from passwords and compare hashes to passwords.
type Hasher interface {
	// GenerateHashFromPassword returns a hash of the password.
	GenerateHashFromPassword(password []byte) (hash []byte, err error)

	// CompareHashWithPassword compares a password hash with a password.
	// Returns nil on success, or an error on failure.
	CompareHashWithPassword(hash, password []byte) (err error)
}
