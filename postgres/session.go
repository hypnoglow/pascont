package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/hypnoglow/pascont/session"
)

const sessionTable = "session"

type sessionRepository struct {
	db *sql.DB
}

// NewSessionRepository returns new session.Repository with PostgreSQL as a storage.
func NewSessionRepository(db *sql.DB) session.Repository {
	return &sessionRepository{db}
}

func (r sessionRepository) Save(s session.Session) error {
	q := fmt.Sprintf(`
		INSERT INTO %s
			(id, account_id, created_at, expires_at)
		VALUES
			($1, $2, $3, $4)
		ON CONFLICT (id) DO
			UPDATE SET
				(account_id, created_at, expires_at)
				= ($2, $3, $4)
			WHERE
				%s.id = $1
	`, pq.QuoteIdentifier(sessionTable), pq.QuoteIdentifier(sessionTable))

	_, err := r.db.Exec(q, s.ID, s.AccountID, s.CreatedAt, s.ExpiresAt)
	return errors.Wrap(err, "Failed to save a session")
}

func (r sessionRepository) FindByID(id string) (*session.Session, error) {
	q := fmt.Sprintf(`
		SELECT
			id, account_id, created_at, expires_at
		FROM
			%s
		WHERE
			id = $1
	`, pq.QuoteIdentifier(sessionTable))

	var s *session.Session
	if err := r.db.QueryRow(q, id).Scan(
		&s.ID,
		&s.AccountID,
		&s.CreatedAt,
		&s.ExpiresAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return s, session.ErrNotFound
		}

		return s, errors.Wrap(err, "Failed to find a session")
	}

	if time.Now().After(s.ExpiresAt) {
		return s, session.ErrExpired
	}

	return s, nil
}

func (r sessionRepository) Delete(id string) error {
	q := fmt.Sprintf(`
		DELETE FROM
			%s
		WHERE
			id = $1
	`, pq.QuoteIdentifier(sessionTable))

	_, err := r.db.Exec(q, id)
	return errors.Wrap(err, "Failed to delete a session")
}

func (r sessionRepository) DeleteAllByAccount(accountID int64) error {
	q := fmt.Sprintf(`
		DELETE FROM
			%s
		WHERE
			account_id = $1
	`, pq.QuoteIdentifier(sessionTable))

	_, err := r.db.Exec(q, accountID)
	return errors.Wrap(err, "Failed to delete account sessions")
}
