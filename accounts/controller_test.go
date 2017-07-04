package accounts

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/hypnoglow/pascont/account"
	"github.com/hypnoglow/pascont/hasher"
)

func TestNewRestController(t *testing.T) {
	NewRestController(
		log.New(ioutil.Discard, "", log.LstdFlags),
		account.NewFakeRepository(nil, nil, nil, nil),
		hasher.NewFakeHasher(nil, nil),
		Options{},
	)
}
