package sessions

import (
	"log"

	"github.com/hypnoglow/pascont/account"
	"github.com/hypnoglow/pascont/hasher"
	"github.com/hypnoglow/pascont/identity"
	"github.com/hypnoglow/pascont/notary"
	"github.com/hypnoglow/pascont/packer"
	"github.com/hypnoglow/pascont/session"
)

const (
	PathSessions = "/sessions"
	PathSession  = "/sessions/"
)

// RestController is a REST controller for sessions.
type RestController struct {
	logger       *log.Logger
	accountRepo  account.Repository
	sessionRepo  session.Repository
	notary       notary.Notary
	packer       packer.Packer
	hasher       hasher.Hasher
	uuidProducer identity.UUIDProducer
	options      Options
}

type Options struct {
	SessionSecretKey []byte
}

// NewRestController returns a new RestController.
func NewRestController(
	logger *log.Logger,
	accountRepo account.Repository,
	sessionRepo session.Repository,
	n notary.Notary,
	p packer.Packer,
	h hasher.Hasher,
	u identity.UUIDProducer,
	opts Options,
) RestController {
	return RestController{
		logger,
		accountRepo,
		sessionRepo,
		n,
		p,
		h,
		u,
		opts,
	}
}
