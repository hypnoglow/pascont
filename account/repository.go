package account

// Repository is a repository for an Account.
type Repository interface {
	// Accept accepts an Application and adds a new Account to the Repository.
	Accept(app Application) (acc *Account, err error)

	// Save saves an Account to the repository.
	// If Account has no ID, returns ErrNoIdentity.
	// Other errors may occur.
	Save(acc Account, passwordHash []byte) error

	// FindWithPasswordHashByUsername retrieves an Account for matching name with it's password.
	// If account with such username not found, returns ErrNotFound.
	// Other errors may occur.
	FindWithPasswordHashByUsername(name string) (account *Account, password []byte, err error)

	// Exists checks whether account with such username exists.
	Exists(username string) (bool, error)
}

// repositoryError is an error occurred in Repository.
type repositoryError string

func (e repositoryError) Error() string {
	return string(e)
}

const (
	// ErrNoIdentity occurs when attempting to save an account with ID equals to 0.
	ErrNoIdentity = repositoryError("Account has no identity")

	// ErrNotFound occurs when account not found.
	ErrNotFound = repositoryError("Account not found")
)
