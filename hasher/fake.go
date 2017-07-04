package hasher

// fakeHasher is a fake Hasher that always produces deterministic given results.
type fakeHasher struct {
	generateHashFromPasswordResults       []FakeGenerateHashFromPasswordResult
	generateHashFromPasswordResultCounter int
	compareHashWithPasswordResults        []FakeCompareHashWithPasswordResult
	compareHashWithPasswordResultCounter  int
}

// FakeGenerateHashFromPasswordResult is a result of hasher.GenerateHashFromPassword.
type FakeGenerateHashFromPasswordResult struct {
	Hash  []byte
	Error error
}

// FakeCompareHashWithPasswordResult is a result of hasher.CompareHashWithPassword.
type FakeCompareHashWithPasswordResult struct {
	Error error
}

// NewFakeHasher returns a new fakeHasher.
func NewFakeHasher(g []FakeGenerateHashFromPasswordResult, c []FakeCompareHashWithPasswordResult) Hasher {
	return &fakeHasher{
		generateHashFromPasswordResults:       g,
		generateHashFromPasswordResultCounter: 0,
		compareHashWithPasswordResults:        c,
		compareHashWithPasswordResultCounter:  0,
	}
}

func (h *fakeHasher) GenerateHashFromPassword(password []byte) (hash []byte, err error) {
	r := h.generateHashFromPasswordResults[h.generateHashFromPasswordResultCounter]
	h.generateHashFromPasswordResultCounter++
	return r.Hash, r.Error
}

func (h *fakeHasher) CompareHashWithPassword(hash, password []byte) (err error) {
	r := h.compareHashWithPasswordResults[h.compareHashWithPasswordResultCounter]
	h.compareHashWithPasswordResultCounter++
	return r.Error
}
