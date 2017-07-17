package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/hypnoglow/pascont/account"
)

const accountTable = "account"

type accountRepository struct {
	db *sql.DB
}

// NewAccountRepository returns a new account.Repository with PostgreSQL as a storage.
func NewAccountRepository(db *sql.DB) account.Repository {
	return &accountRepository{db: db}
}

func (r accountRepository) Accept(app account.Application) (acc *account.Account, err error) {
	q := fmt.Sprintf(`
		INSERT INTO %s
			(name, password_hash, created_at, updated_at)
		VALUES
			($1, $2, $3, $3)
		RETURNING id, name, created_at, updated_At
	`, pq.QuoteIdentifier(accountTable))

	err = r.db.QueryRow(q, app.Name, app.PasswordHash, app.CreatedAt).Scan(
		&acc.ID,
		&acc.Name,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)

	return acc, errors.Wrap(err, "Failed to accept an application")
}

func (r accountRepository) Save(acc account.Account, passwordHash []byte) error {
	if acc.ID == 0 {
		return account.ErrNoIdentity
	}

	q := fmt.Sprintf(`
		UPDATE %s
		SET (name, password_hash, created_at, updated_at)
			= ($2, $3, $4, $5)
		WHERE
			%s.id = $1
		RETURNING id

	`, pq.QuoteIdentifier(accountTable), pq.QuoteIdentifier(accountTable))

	_, err := r.db.Exec(q, acc.ID, acc.Name, passwordHash, acc.CreatedAt, acc.UpdatedAt)
	return errors.Wrap(err, "Failed to save an account")
}

func (r accountRepository) FindWithPasswordHashByUsername(name string) (*account.Account, []byte, error) {
	q := fmt.Sprintf(`
		SELECT
			id, name, password_hash, created_at, updated_at
		FROM
			%s
		WHERE
			name = $1

	`, pq.QuoteIdentifier(accountTable))

	acc := &account.Account{}
	var passwordHash []byte

	err := r.db.QueryRow(q, name).Scan(
		&acc.ID,
		&acc.Name,
		&passwordHash,
		&acc.CreatedAt,
		&acc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		err = account.ErrNotFound
	} else {
		err = errors.Wrap(err, "Failed to find an account")
	}

	return acc, passwordHash, err
}

func (r accountRepository) Exists(name string) (bool, error) {
	q := fmt.Sprintf(`
		SELECT
			id
		FROM
			%s
		WHERE
			name = $1
	`, pq.QuoteIdentifier(accountTable))

	var id int64
	if err := r.db.QueryRow(q, name).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, errors.Wrap(err, "Failed to check if account exists")
	}

	return true, nil
}
