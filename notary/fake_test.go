package notary

import (
	"bytes"
	"testing"
)

func TestFakeNotary(t *testing.T) {
	key := []byte("key")
	message := []byte("message")

	cases := []struct {
		signResult   FakeNotarySignResult
		verifyResult FakeNotaryVerifyResult
	}{
		{
			signResult: FakeNotarySignResult{
				Signature: []byte("signature"),
			},
			verifyResult: FakeNotaryVerifyResult{
				Result: true,
			},
		},
		{
			signResult: FakeNotarySignResult{
				Signature: []byte("other signature"),
			},
			verifyResult: FakeNotaryVerifyResult{
				Result: false,
			},
		},
	}

	for i, c := range cases {
		n := NewFakeNotary(
			[]FakeNotarySignResult{c.signResult},
			[]FakeNotaryVerifyResult{c.verifyResult},
		)

		signature := n.Sign(message, key)
		if !bytes.Equal(signature, c.signResult.Signature) {
			t.Errorf(
				"testcase %d: Expected n.Sign() to equal %v, but got %v\n",
				i,
				c.signResult.Signature,
				signature,
			)
		}

		verified := n.Verify(message, signature, key)
		if verified != c.verifyResult.Result {
			t.Errorf(
				"testcase %d: Expected n.Verify() to equal %v, but got %v\n",
				i,
				c.verifyResult.Result,
				verified,
			)
		}
	}
}
