package account

import "time"

// Application represents a new account proposal.
// It can be used as intention to create a new Account.
type Application struct {
	Name         string
	PasswordHash []byte
	CreatedAt    time.Time
}

// NewApplication returns a new Application.
func NewApplication(name string, passwordHash []byte, createdAt time.Time) Application {
	createdAt = createdAt.UTC().Truncate(time.Second)
	return Application{name, passwordHash, createdAt}
}
