package packer

type fakePacker struct {
	packResults         []FakePackerPackResult
	packResultCounter   int
	unpackResults       []FakePackerUnpackResult
	unpackResultCounter int
}

type FakePackerPackResult struct {
	Pack  []byte
	Error error
}

type FakePackerUnpackResult struct {
	Message   []byte
	Signature []byte
	Error     error
}

// NewFakePacker returns a new fake Packer.
func NewFakePacker(packResults []FakePackerPackResult, unpackResults []FakePackerUnpackResult) Packer {
	return &fakePacker{
		packResults:         packResults,
		packResultCounter:   0,
		unpackResults:       unpackResults,
		unpackResultCounter: 0,
	}
}

func (p *fakePacker) Pack(message, signature []byte) (pack []byte, err error) {
	res := p.packResults[p.packResultCounter]
	p.packResultCounter++
	return res.Pack, res.Error
}

func (p *fakePacker) Unpack(pack []byte) (message, signature []byte, err error) {
	res := p.unpackResults[p.unpackResultCounter]
	p.unpackResultCounter++
	return res.Message, res.Signature, res.Error
}
