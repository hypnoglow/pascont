package sessions

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/hypnoglow/pascont/account"
	"github.com/hypnoglow/pascont/hasher"
	"github.com/hypnoglow/pascont/identity"
	"github.com/hypnoglow/pascont/notary"
	"github.com/hypnoglow/pascont/packer"
	"github.com/hypnoglow/pascont/session"
)

func TestNewRestController(t *testing.T) {
	NewRestController(
		log.New(ioutil.Discard, "", log.LstdFlags),
		account.NewFakeRepository(nil, nil, nil, nil),
		session.NewFakeRepository(nil, nil),
		notary.NewFakeNotary(nil, nil),
		packer.NewFakePacker(nil, nil),
		hasher.NewFakeHasher(nil, nil),
		identity.NewFakeIdentificatory(nil),
		Options{
			SessionSecretKey: []byte("secret_key"),
		},
	)
}
