package notary

// Notary can sign message with a key and verify if message is signed and can be trusted.
type Notary interface {
	// Sign returns signature for the message.
	Sign(message, key []byte) (signature []byte)

	// Verify verifies the message by the signature.
	Verify(message, messageMAC, key []byte) bool
}
