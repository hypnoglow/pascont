package account

import "time"

// Account represents account data.
type Account struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAccount returns a new Account.
func NewAccount(id int64, name string, createdAt, updatedAt time.Time) *Account {
	createdAt = createdAt.UTC().Truncate(time.Second)
	updatedAt = updatedAt.UTC().Truncate(time.Second)
	return &Account{id, name, createdAt, updatedAt}
}
