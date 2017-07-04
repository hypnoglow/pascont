package session

// Repository is a repository for Session.
type Repository interface {
	// Save saves a session to the repository.
	Save(sess Session) error

	// FindByID looks for a session with specified id.
	// If session not found, returns ErrNotFound.
	// If session is already expired, returns ErrExpired.
	FindByID(id string) (*Session, error)

	// Delete removes the session.
	Delete(id string) error

	// DeleteAllByAccount removes all sessions of the account.
	DeleteAllByAccount(accountID int64) error
}

// repositoryError is an error occured in Repository
type repositoryError string

func (e repositoryError) Error() string {
	return string(e)
}

const (
	// ErrNotFound occurs when session not found.
	ErrNotFound = repositoryError("Session not found")

	// ErrExpired occurs when session exists in repository but expired.
	ErrExpired = repositoryError("Session is expired")
)
