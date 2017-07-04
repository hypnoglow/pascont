package packer

import (
	"bytes"
	"testing"
)

func TestBase64Packer(t *testing.T) {
	cases := []struct {
		mlen      int
		message   []byte
		signature []byte
	}{
		{
			mlen:      5,
			message:   []byte("hello"),
			signature: []byte("signature"),
		},
	}

	for i, c := range cases {
		p := NewBase64Packer(c.mlen)

		pack, err := p.Pack(c.message, c.signature)
		if err != nil {
			t.Errorf("testcase %c: Expected err to be nil, got %v\n", i, err)
		}

		m, s, err := p.Unpack(pack)
		if err != nil {
			t.Errorf("testcase %c: Expected err to be nil, got %v\n", i, err)
		}

		if !bytes.Equal(m, c.message) {
			t.Errorf(
				"testcase %c: Expected message %v to be same after Pack-Unpack cycle, but got %v\n",
				i,
				c.message,
				m,
			)
		}

		if !bytes.Equal(s, c.signature) {
			t.Errorf(
				"testcase %c: Expected signature %v to be same after Pack-Unpack cycle, but got %v\n",
				i,
				c.signature,
				s,
			)
		}
	}
}

func TestBase64Packer_Pack(t *testing.T) {
	// Test that Pack fails on wrong message length.

	p := NewBase64Packer(4)
	_, err := p.Pack([]byte("Hello"), []byte("Signature"))
	if err == nil {
		t.Errorf("Expected Pack to fail due to wrong message length")
	}
}

func TestBase64Packer_Unpack(t *testing.T) {
	// Test that Unpack fails on invalid base64 data.

	p := NewBase64Packer(4)
	_, _, err := p.Unpack([]byte("Invalid base64 data"))
	if err == nil {
		t.Errorf("Expected Unpack to fail due to invalid base64 data")
	}
}
