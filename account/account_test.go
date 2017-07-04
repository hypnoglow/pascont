package account

import (
	"testing"
	"time"
)

func TestNewAccount(t *testing.T) {
	now := time.Now()

	cases := []struct {
		id                int64
		name              string
		createdAt         time.Time
		updatedAt         time.Time
		expectedCreatedAt time.Time
		expectedUpdatedAt time.Time
	}{
		{
			id:                123,
			name:              "email@email.com",
			createdAt:         now,
			updatedAt:         now,
			expectedCreatedAt: now.UTC().Truncate(time.Second),
			expectedUpdatedAt: now.UTC().Truncate(time.Second),
		},
	}

	for i, c := range cases {
		actual := NewAccount(c.id, c.name, c.createdAt, c.updatedAt)

		if actual.ID != c.id {
			t.Errorf(
				"testcase %d: Expected ID to be %v but got %v\n",
				i,
				c.id,
				actual.ID,
			)
		}

		if actual.Name != c.name {
			t.Errorf(
				"testcase %d: Expected Name to be %v but got %v\n",
				i,
				c.name,
				actual.Name,
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

		if actual.UpdatedAt != c.expectedUpdatedAt {
			t.Errorf(
				"testcase %d: Expected UpdatedAt to be %v but got %v\n",
				i,
				c.expectedUpdatedAt,
				actual.UpdatedAt,
			)
		}
	}
}
