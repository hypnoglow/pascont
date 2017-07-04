package packer

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestFakePacker(t *testing.T) {
	cases := []struct {
		packResult   FakePackerPackResult
		unpackResult FakePackerUnpackResult
	}{
		{
			packResult: FakePackerPackResult{
				Pack:  []byte("pack"),
				Error: nil,
			},
			unpackResult: FakePackerUnpackResult{
				Message:   []byte("message"),
				Signature: []byte("signature"),
				Error:     nil,
			},
		},
		{
			packResult: FakePackerPackResult{
				Pack:  []byte("test pack"),
				Error: fmt.Errorf("Pack error"),
			},
			unpackResult: FakePackerUnpackResult{
				Message:   []byte("test message"),
				Signature: []byte("test signature"),
				Error:     fmt.Errorf("Unpack error"),
			},
		},
	}

	for i, c := range cases {
		p := NewFakePacker(
			[]FakePackerPackResult{c.packResult},
			[]FakePackerUnpackResult{c.unpackResult},
		)

		pack, err := p.Pack([]byte("message"), []byte("signature"))

		if !bytes.Equal(pack, c.packResult.Pack) {
			t.Errorf(
				"testcase %d: Expected packer.Pack() to equal %v, but got %v\n",
				i,
				c.packResult.Pack,
				pack,
			)
		}

		if !reflect.DeepEqual(err, c.packResult.Error) {
			t.Errorf(
				"testcase %d: Expected packer.Pack() error to equal %v, but got %v\n",
				i,
				c.packResult.Error,
				err,
			)
		}

		mes, sig, err := p.Unpack(pack)

		if !bytes.Equal(mes, c.unpackResult.Message) {
			t.Errorf(
				"testcase %d: Expected packer.Unpack() message to equal %v, but got %v\n",
				i,
				c.unpackResult.Message,
				mes,
			)
		}

		if !bytes.Equal(sig, c.unpackResult.Signature) {
			t.Errorf(
				"testcase %d: Expected packer.Unpack() signature to equal %v, but got %v\n",
				i,
				c.unpackResult.Signature,
				sig,
			)
		}

		if !reflect.DeepEqual(err, c.unpackResult.Error) {
			t.Errorf(
				"testcase %d: Expected packer.Unpack() error to equal %v, but got %v\n",
				i,
				c.unpackResult.Error,
				err,
			)
		}
	}
}
