package packer

// Packer can encode and decode message and it's signature.
type Packer interface {
	// Pack encodes message and it's signature.
	Pack(message, signature []byte) (pack []byte, err error)

	// Unpack decodes pack to message and it's signature
	Unpack(pack []byte) (message, signature []byte, err error)
}
