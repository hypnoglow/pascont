package account

type fakeRepository struct {
	acceptResults                               []FakeRepositoryAcceptResult
	acceptResultCounter                         int
	saveResults                                 []FakeRepositorySaveResult
	saveResultCounter                           int
	existsResults                               []FakeRepositoryExistsResult
	existsResultCounter                         int
	findWithPasswordHashByUsernameResults       []FakeRepositoryFindWithPasswordHashByUsernameResult
	findWithPasswordHashByUsernameResultCounter int
}

type FakeRepositoryAcceptResult struct {
	Account *Account
	Error   error
}

type FakeRepositorySaveResult struct {
	Error error
}

type FakeRepositoryExistsResult struct {
	Exists bool
	Error  error
}

type FakeRepositoryFindWithPasswordHashByUsernameResult struct {
	Account      *Account
	PasswordHash []byte
	Error        error
}

// NewFakeRepository returns a new fake Repository.
func NewFakeRepository(
	acceptResults []FakeRepositoryAcceptResult,
	saveResults []FakeRepositorySaveResult,
	existsResults []FakeRepositoryExistsResult,
	findWithPasswordHashByUsernameResults []FakeRepositoryFindWithPasswordHashByUsernameResult,
) Repository {
	return &fakeRepository{
		acceptResults:                               acceptResults,
		acceptResultCounter:                         0,
		saveResults:                                 saveResults,
		saveResultCounter:                           0,
		existsResults:                               existsResults,
		existsResultCounter:                         0,
		findWithPasswordHashByUsernameResults:       findWithPasswordHashByUsernameResults,
		findWithPasswordHashByUsernameResultCounter: 0,
	}
}

func (r *fakeRepository) Accept(app Application) (acc *Account, err error) {
	res := r.acceptResults[r.acceptResultCounter]
	r.acceptResultCounter++
	return res.Account, res.Error
}

func (r *fakeRepository) Save(acc Account, passwordHash []byte) error {
	res := r.saveResults[r.saveResultCounter]
	r.saveResultCounter++
	return res.Error
}

func (r *fakeRepository) FindWithPasswordHashByUsername(name string) (*Account, []byte, error) {
	res := r.findWithPasswordHashByUsernameResults[r.findWithPasswordHashByUsernameResultCounter]
	r.findWithPasswordHashByUsernameResultCounter++
	return res.Account, res.PasswordHash, res.Error
}

func (r *fakeRepository) Exists(username string) (bool, error) {
	res := r.existsResults[r.existsResultCounter]
	r.existsResultCounter++
	return res.Exists, res.Error
}
