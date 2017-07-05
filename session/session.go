package session

import (
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/hypnoglow/pascont/identity"
	"github.com/hypnoglow/pascont/notary"
	"github.com/hypnoglow/pascont/packer"
)

const (
	// SessionIDLength is the length of a session ID in bytes.
	SessionIDLength = 36

	// SessionExpiresAtLength is the length of a session ExpiresAt in Unix Timestamp form in bytes.
	SessionExpiresAtLength = 10
)

// Session represents session data.
type Session struct {
	ID        string
	AccountID int64
	CreatedAt time.Time
	ExpiresAt time.Time
}

// NewSession returns a new Session.
func NewSession(id string, accountID int64, createdAt, expiresAt time.Time) *Session {
	createdAt = createdAt.UTC().Truncate(time.Second)
	expiresAt = expiresAt.UTC().Truncate(time.Second)
	return &Session{
		ID:        id,
		AccountID: accountID,
		CreatedAt: createdAt,
		ExpiresAt: expiresAt,
	}
}

// Token returns the token that represents the Session.
// Token includes Session ID, ExpiresAt and a signature.
func (s Session) Token(notary notary.Notary, packer packer.Packer, secretKey []byte) (token string, err error) {
	timestamp := make([]byte, 10)
	copy(timestamp, []byte(strconv.FormatInt(s.ExpiresAt.Unix(), 10)))
	message := append([]byte(s.ID), timestamp...)

	signature := notary.Sign(message, secretKey)
	pack, err := packer.Pack(message, signature)
	if err != nil {
		return "", errors.Wrap(err, "Failed to pack message with token and it's signature")
	}

	return string(pack), nil
}

// ResetExpiresAt resets session expire date to current time plus duration.
func (s *Session) ResetExpiresAt(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration).UTC()
}

// CreateSession creates a new Session.
// The Session ID is a UUID produced by identity.UUIDProducer func.
// The CreatedAt is set to the current time.
// The ExpiresAt is set to the current time plus duration.
func CreateSession(uuidProducer identity.UUIDProducer, accountID int64, duration time.Duration) *Session {
	now := time.Now()
	return NewSession(
		uuidProducer(),
		accountID,
		now,
		now.Add(duration),
	)
}
