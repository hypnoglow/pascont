package accounts

import (
	"log"

	"github.com/hypnoglow/pascont/account"
	"github.com/hypnoglow/pascont/hasher"
)

const (
	PathAccounts = "/accounts"
)

// RestController is a REST controller for accounts.
type RestController struct {
	errorLogger *log.Logger
	accountRepo account.Repository
	hasher      hasher.Hasher
	options     Options
}

// Options is a structure holding accounts RestController specific options.
type Options struct {
	// no options currently
}

// NewRestController returns a new RestController.
func NewRestController(
	logger *log.Logger,
	accountRepo account.Repository,
	hasher hasher.Hasher,
	opts Options,
) RestController {
	return RestController{logger, accountRepo, hasher, opts}
}
