package notary

type fakeNotary struct {
	signResults         []FakeNotarySignResult
	signResultCounter   int
	verifyResults       []FakeNotaryVerifyResult
	verifyResultCounter int
}

type FakeNotarySignResult struct {
	Signature []byte
}

type FakeNotaryVerifyResult struct {
	Result bool
}

// NewFakeNotary returns a new fake Notary.
func NewFakeNotary(signResults []FakeNotarySignResult, verifyResults []FakeNotaryVerifyResult) Notary {
	return &fakeNotary{
		signResults:         signResults,
		signResultCounter:   0,
		verifyResults:       verifyResults,
		verifyResultCounter: 0,
	}
}

func (n *fakeNotary) Sign(message, key []byte) (signature []byte) {
	res := n.signResults[n.signResultCounter]
	n.signResultCounter++
	return res.Signature
}

func (n *fakeNotary) Verify(message, messageMAC, key []byte) bool {
	res := n.verifyResults[n.verifyResultCounter]
	n.verifyResultCounter++
	return res.Result
}
