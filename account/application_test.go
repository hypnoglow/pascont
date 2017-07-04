package account

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestNewApplication(t *testing.T) {
	now := time.Now()

	cases := []struct {
		name              string
		passwordHash      []byte
		createdAt         time.Time
		expectedCreatedAt time.Time
	}{
		{
			name:              "email@email.com",
			passwordHash:      []byte("some_secure_hash"),
			createdAt:         now,
			expectedCreatedAt: now.UTC().Truncate(time.Second),
		},
	}

	for i, c := range cases {
		actual := NewApplication(c.name, c.passwordHash, c.createdAt)

		if actual.Name != c.name {
			t.Errorf(
				"testcase %d: Expected Name to be %v but got %v\n",
				i,
				c.name,
				actual.Name,
			)
		}

		if !bytes.Equal(actual.PasswordHash, c.passwordHash) {
			t.Errorf(
				"testcase %d: Expected PasswordHash to be %v but got %v\n",
				i,
				c.passwordHash,
				actual.PasswordHash,
			)
		}

		if !reflect.DeepEqual(actual.CreatedAt, c.expectedCreatedAt) {
			t.Errorf(
				"testcase %d: Expected CreatedAt to be %v, but got %v\n",
				i,
				c.expectedCreatedAt,
				actual.CreatedAt,
			)
		}
	}
}
