package identity

import "github.com/pborman/uuid"

type UUIDProducer func() string

func NewUUIDV4() string {
	return uuid.New()
}
