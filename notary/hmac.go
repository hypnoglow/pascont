package notary

import (
	"crypto/hmac"
	"crypto/sha512"
)

type hmacNotary struct{}

// NewHMACNotary returnw new Notary which uses HMAC to sign and verify messages.
func NewHMACNotary() Notary {
	return hmacNotary{}
}

func (n hmacNotary) Sign(message, key []byte) (signature []byte) {
	mac := hmac.New(sha512.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}

func (n hmacNotary) Verify(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha512.New, key)
	mac.Write(message)
	return hmac.Equal(messageMAC, mac.Sum(nil))
}
