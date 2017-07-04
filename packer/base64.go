package packer

import (
	"encoding/base64"
	"fmt"
)

type base64Packer struct {
	mlen int
}

// NewBase64Packer returns new Packer which packs using base64 encoding.
func NewBase64Packer(mlen int) Packer {
	return base64Packer{mlen}
}

func (p base64Packer) Pack(message, signature []byte) (pack []byte, err error) {
	if len(message) != p.mlen {
		return nil, fmt.Errorf("message len is not equal to packer mlen")
	}

	m := append(message, signature...)
	pack = make([]byte, base64.URLEncoding.EncodedLen(len(m)))
	base64.URLEncoding.Encode(pack, m)
	return pack, nil
}

func (p base64Packer) Unpack(pack []byte) (message, signature []byte, err error) {
	decoded := make([]byte, base64.URLEncoding.DecodedLen(len(pack)))
	n, err := base64.URLEncoding.Decode(decoded, pack)
	if err != nil {
		return nil, nil, err
	}

	return decoded[:p.mlen], decoded[p.mlen:n], nil
}
