package identity

// Identificatory can issue new UUIDs.
type Identificatory interface {
	// NewUUID returns a new uuid in string representation.
	NewUUID() string
}
