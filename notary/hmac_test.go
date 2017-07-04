package notary

import "testing"

func TestHMACNotary(t *testing.T) {
	cases := []struct {
		key     []byte
		message []byte
	}{
		{
			[]byte("123"),
			[]byte("Hello"),
		},
		{
			[]byte("secret_key"),
			[]byte("Message to sign"),
		},
	}

	n := NewHMACNotary()

	for i, c := range cases {
		signature := n.Sign(c.message, c.key)
		if !n.Verify(c.message, signature, c.key) {
			t.Errorf(
				"testcase %c: Expected message to be signed and verified, key: %s, message: %s, signature: %s",
				i,
				c.key,
				c.message,
				signature,
			)
		}
	}
}
