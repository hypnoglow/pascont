package identity

type fakeIdentificatory struct {
	newUUIDResults       []FakeNewUUIDResult
	newUUIDResultCounter int
}

type FakeNewUUIDResult struct {
	UUID string
}

// NewFakeIdentificatory retruns new fake Identificatory.
func NewFakeIdentificatory(newUUIDResults []FakeNewUUIDResult) Identificatory {
	return &fakeIdentificatory{
		newUUIDResults:       newUUIDResults,
		newUUIDResultCounter: 0,
	}
}

func (f *fakeIdentificatory) NewUUID() string {
	res := f.newUUIDResults[f.newUUIDResultCounter]
	f.newUUIDResultCounter++
	return res.UUID
}
