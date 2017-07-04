package sessions

import "time"

const (
	// sessionDefaultDuration is a default duration of a session.
	sessionDefaultDuration = time.Hour * 24 * 3

	// sessionMaxDuration is a max duration of a session.
	sessionMaxDuration = time.Hour * 24 * 3
)
