package session

import (
	"fmt"
	"testing"
	"time"

	"github.com/hypnoglow/pascont/identity"
	"github.com/hypnoglow/pascont/notary"
	"github.com/hypnoglow/pascont/packer"
	"github.com/pkg/errors"
)

func TestNewSession(t *testing.T) {
	now := time.Now()
	duration := time.Hour

	cases := []struct {
		id                string
		accountID         int64
		createdAt         time.Time
		expiresAt         time.Time
		expectedCreatedAt time.Time
		expectedExpiresAt time.Time
	}{
		{
			id:                "123-456-789",
			accountID:         123,
			createdAt:         now,
			expiresAt:         now.Add(duration),
			expectedCreatedAt: now.UTC().Truncate(time.Second),
			expectedExpiresAt: now.Add(duration).UTC().Truncate(time.Second),
		},
	}

	for i, c := range cases {
		actual := NewSession(c.id, c.accountID, c.createdAt, c.expiresAt)

		if actual.ID != c.id {
			t.Errorf(
				"testcase %d: Expected ID to be %v but got %v\n",
				i,
				c.id,
				actual.ID,
			)
		}

		if actual.AccountID != c.accountID {
			t.Errorf(
				"testcase %d: Expected AccountID to be %v but got %v\n",
				i,
				c.id,
				actual.ID,
			)
		}

		if actual.CreatedAt != c.expectedCreatedAt {
			t.Errorf(
				"testcase %d: Expected CreatedAt to be %v but got %v\n",
				i,
				c.expectedCreatedAt,
				actual.CreatedAt,
			)
		}

		if actual.ExpiresAt != c.expectedExpiresAt {
			t.Errorf(
				"testcase %d: Expected ExpiresAt to be %v but got %v\n",
				i,
				c.expectedExpiresAt,
				actual.ExpiresAt,
			)
		}
	}
}

func TestSession_Token(t *testing.T) {
	n := notary.NewHMACNotary()                                           // TODO: replace with fake
	p := packer.NewBase64Packer(SessionIDLength + SessionExpiresAtLength) // TODO: replace with fake
	key := []byte("secret_key")
	now := time.Now().UTC().Truncate(time.Second)

	cases := []struct {
		sess        *Session
		expectedErr error
	}{
		{
			sess:        NewSession("12345678-90ab-cdef-0123-4567890abcde", 123, now, now.Add(time.Hour)),
			expectedErr: nil,
		},
		// Session with incorrect ID.
		{
			sess: NewSession("sessID", 123, time.Now(), time.Now()),
			expectedErr: errors.Wrap(
				fmt.Errorf("message len is not equal to packer mlen"),
				"Failed to pack message with token and it's signature",
			),
		},
	}

	for i, c := range cases {
		_, err := c.sess.Token(n, p, key)
		if err == nil && c.expectedErr != nil || err != nil && err.Error() != c.expectedErr.Error() {
			t.Errorf(
				"testcase %d: Expected err to be %v, but got %v\n",
				i,
				c.expectedErr,
				err,
			)
		}
	}
}

func TestSession_ResetExpiresAt(t *testing.T) {
	extendDuration := time.Second * 60
	now := time.Now().UTC().Truncate(time.Second)

	s := NewSession("12345678-90ab-cdef-0123-4567890abcde", 123, now, now.Add(time.Hour))
	s.ResetExpiresAt(extendDuration)

	// ExpiresAt should be time.Now() + extendDuration
	expected := time.Now().UTC().Add(extendDuration)
	if expected.Sub(s.ExpiresAt) > time.Second {
		t.Errorf(
			"Expected session to have ExpiresAt %v, but got %v\n",
			expected,
			s.ExpiresAt,
		)
	}
}

func TestCreateSession(t *testing.T) {
	cases := []struct {
		idf       identity.Identificatory
		accountID int64
		duration  time.Duration
		expected  *Session
	}{
		{
			idf: identity.NewFakeIdentificatory(
				[]identity.FakeNewUUIDResult{
					{
						UUID: "12345678-90ab-cdef-0123-4567890abcde",
					},
				},
			),
			accountID: 123,
			duration:  time.Second * 30,
			expected: &Session{
				ID:        "12345678-90ab-cdef-0123-4567890abcde",
				AccountID: 123,
				CreatedAt: time.Now().UTC().Truncate(time.Second),
				ExpiresAt: time.Now().Add(time.Second * 30).UTC().Truncate(time.Second),
			},
		},
	}

	for i, c := range cases {
		actual := CreateSession(c.idf.NewUUID, c.accountID, c.duration)

		if actual.ID != c.expected.ID {
			t.Errorf(
				"testcase %d: Expected session to have ID %v but got %v\n",
				i,
				c.expected.ID,
				actual.ID,
			)
		}

		if actual.AccountID != c.expected.AccountID {
			t.Errorf(
				"testcase %d: Expected session to have AccountID %v, but got %v\n",
				i,
				c.expected.AccountID,
				actual.AccountID,
			)
		}

		if actual.CreatedAt.Sub(c.expected.CreatedAt) > time.Second {
			t.Errorf(
				"testcase %d: Expected session to have CreatedAt %v, but got %v\n",
				i,
				c.expected.CreatedAt,
				actual.CreatedAt,
			)
		}

		if actual.ExpiresAt.Sub(c.expected.ExpiresAt) > time.Second {
			t.Errorf(
				"testcase %d: Expected session to have ExpiresAt %v, but got %v\n",
				i,
				c.expected.ExpiresAt,
				actual.ExpiresAt,
			)
		}
	}
}
