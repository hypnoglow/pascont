package identity

import "github.com/pborman/uuid"

type uuidV4 struct{}

// NewUUIDV4 returns a new Identificatory which uses uuidv4 to issue uuids.
func NewUUIDV4() Identificatory {
	return uuidV4{}
}

func (u uuidV4) NewUUID() string {
	return uuid.New()
}
