package kit

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// NewJSONBody encodes and returns payload as json body.
func NewJSONBody(payload interface{}) (body []byte) {
	if payload == nil {
		return nil
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(errors.Wrap(err, "Failed to encode body to JSON"))
	}

	return body
}
