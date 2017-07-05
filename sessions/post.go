package sessions

import (
	"net/http"
	"time"

	"github.com/hypnoglow/pascont/account"
	"github.com/hypnoglow/pascont/kit"
	"github.com/hypnoglow/pascont/kit/form"
	"github.com/hypnoglow/pascont/kit/schema"
	"github.com/hypnoglow/pascont/session"
)

// PostSessions is a handler for:
// POST /sessions
func (c RestController) PostSessions(w http.ResponseWriter, req *http.Request) {
	var sessForm postSessionForm
	form.PopulateFormFromJSON(req.Body, &sessForm)
	if !sessForm.Validate() {
		kit.RespondWithFormErrors(w, http.StatusBadRequest, sessForm.ValidationErrors())
		return
	}

	// Name+PasswordHash pair must be correct.
	acc, passwordHash, err := c.accountRepo.FindWithPasswordHashByUsername(sessForm.Name)
	if err != nil {
		if err == account.ErrNotFound {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			c.logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if err := c.hasher.CompareHashWithPassword(passwordHash, []byte(sessForm.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sess := session.CreateSession(c.uuidProducer, acc.ID, sessionDefaultDuration)
	if err := c.sessionRepo.Save(*sess); err != nil {
		c.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := sess.Token(c.notary, c.packer, c.options.SessionSecretKey)
	if err != nil {
		c.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	kit.RespondJSON(w, http.StatusCreated, schema.NewResultBody(
		postSessionSchema{
			Token:     token,
			ID:        sess.ID,
			AccountID: sess.AccountID,
			CreatedAt: sess.CreatedAt,
			ExpiresAt: sess.ExpiresAt,
		},
		nil,
	))
}

type postSessionForm struct {
	form.BaseForm
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (f *postSessionForm) Validate() bool {
	if len(f.Name) == 0 {
		f.AddError("Name must not be empty", "name", f.Name)
	}

	if len(f.Password) == 0 {
		f.AddError("PasswordHash must not be empty", "password", "")
	}

	return len(f.ValidationErrors()) == 0
}

type postSessionSchema struct {
	Token     string    `json:"token"`
	ID        string    `json:"id"`
	AccountID int64     `json:"accountID"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
