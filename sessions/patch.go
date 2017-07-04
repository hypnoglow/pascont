package sessions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hypnoglow/pascont/kit"
	"github.com/hypnoglow/pascont/kit/form"
	"github.com/hypnoglow/pascont/kit/middleware"
	"github.com/hypnoglow/pascont/kit/schema"
	"github.com/hypnoglow/pascont/session"
)

// PatchSession is a handler for:
// PATCH /session/:id
func (c RestController) PatchSession(w http.ResponseWriter, req *http.Request) {
	// Check that it is actually a UUID.
	sid, ok := req.Context().Value(middleware.ContextKeyPathID{}).(string)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sidFromToken, ok := req.Context().Value(middleware.ContextKeySessionID{}).(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Session ID from path and from auth token MUST match in order to proceed.
	if sidFromToken != sid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var sessForm patchSessionForm
	form.PopulateFormFromJSON(req.Body, &sessForm)
	if !sessForm.Validate() {
		kit.RespondWithFormErrors(w, http.StatusBadRequest, sessForm.ValidationErrors())
		return
	}

	sess, err := c.sessionRepo.FindByID(sid)
	if err != nil {
		if err == session.ErrNotFound || err == session.ErrExpired {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		c.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sess.ExpiresAt = sessForm.ExpiresAt.UTC()

	if err := c.sessionRepo.Save(*sess); err != nil {
		c.logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Make new token, because ExpiresAt changed.
	newToken, err := sess.Token(c.notary, c.packer, c.options.SessionSecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	kit.RespondJSON(w, http.StatusOK, schema.NewResultBody(
		patchSessionSchema{
			Token:     newToken,
			ID:        sess.ID,
			AccountID: sess.AccountID,
			CreatedAt: sess.CreatedAt,
			ExpiresAt: sess.ExpiresAt,
		},
		nil,
	))
}

type patchSessionForm struct {
	form.BaseForm
	ExpiresAt time.Time `json:"expiresAt"`
}

func (f *patchSessionForm) Validate() bool {
	if time.Now().After(f.ExpiresAt) {
		f.AddError(
			fmt.Sprintf(
				"ExpiresAt must be greater than %s",
				time.Now().UTC().Format(time.RFC3339),
			),
			"expiresAt",
			f.ExpiresAt,
		)
	}

	if f.ExpiresAt.Sub(time.Now()) > sessionMaxDuration {
		f.AddError(
			fmt.Sprintf(
				"ExpiresAt must be not greater than %s",
				time.Now().UTC().Add(sessionMaxDuration).Format(time.RFC3339),
			),
			"expiresAt",
			f.ExpiresAt,
		)
	}

	return len(f.ValidationErrors()) == 0
}

type patchSessionSchema struct {
	Token     string    `json:"token"`
	ID        string    `json:"id"`
	AccountID int64     `json:"accountID"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
