package session

import (
	"fmt"

	"github.com/hypnoglow/pascont/notary"
	"github.com/hypnoglow/pascont/packer"
)

// ExtractIDFromToken extracts a session ID from a token string.
// It unpacks the token and checks the signature.
// If the token is invalid, returns an error.
func ExtractIDFromToken(token string, p packer.Packer, n notary.Notary, secretKey []byte) (id string, err error) {
	sid, signature, err := p.Unpack([]byte(token))
	if err != nil {
		return "", fmt.Errorf("Failed to decode token")
	}

	verified := n.Verify(sid, signature, secretKey)
	if !verified {
		return "", fmt.Errorf("Failed to verify session ID")
	}

	return string(sid[:SessionIDLength]), nil
}

// TokenExtractor returns a func which can be used to extract a session ID from a token.
func TokenExtractor(p packer.Packer, n notary.Notary, key []byte) func(token string) (id string, err error) {
	return func(token string) (id string, err error) {
		return ExtractIDFromToken(token, p, n, key)
	}
}
