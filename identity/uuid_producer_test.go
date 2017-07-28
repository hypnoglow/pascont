package identity

import (
	"testing"

	"github.com/pborman/uuid"
)

func TestNewUUIDV4(t *testing.T) {
	if uuid.Parse(NewUUIDV4()) == nil {
		t.Errorf("UUID is incorrect")
	}
}