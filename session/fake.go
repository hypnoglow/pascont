package session

type fakeRepository struct {
	saveResults           []FakeRepositorySaveResult
	saveResultCounter     int
	findByIDResults       []FakeRepositoryFindByIDResult
	findByIDResultCounter int
}

type FakeRepositorySaveResult struct {
	Error error
}

type FakeRepositoryFindByIDResult struct {
	Session *Session
	Error   error
}

// NewFakeRepository returns a new fake Repository.
func NewFakeRepository(
	saveResults []FakeRepositorySaveResult,
	findByIDResults []FakeRepositoryFindByIDResult,
) Repository {
	return &fakeRepository{
		saveResults:           saveResults,
		saveResultCounter:     0,
		findByIDResults:       findByIDResults,
		findByIDResultCounter: 0,
	}
}

func (r *fakeRepository) Save(sess Session) error {
	res := r.saveResults[r.saveResultCounter]
	r.saveResultCounter++
	return res.Error
}

func (r *fakeRepository) FindByID(id string) (*Session, error) {
	res := r.findByIDResults[r.findByIDResultCounter]
	r.findByIDResultCounter++
	return res.Session, res.Error
}

func (r *fakeRepository) Delete(id string) error {
	// NOT FAKED YET
	return nil
}

func (r *fakeRepository) DeleteAllByAccount(accountID int64) error {
	// NOT FAKED YET
	return nil
}
