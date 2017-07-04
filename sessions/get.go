package sessions

import (
	"net/http"
	"strings"
	"time"

	"github.com/hypnoglow/pascont/kit"
	"github.com/hypnoglow/pascont/kit/middleware"
	"github.com/hypnoglow/pascont/kit/schema"
	"github.com/hypnoglow/pascont/session"
)

// GetSession is a handler for:
// GET /session/:id
func (c RestController) GetSession(w http.ResponseWriter, req *http.Request) {
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

	kit.RespondJSON(w, http.StatusOK, schema.NewResultBody(
		getSessionSchema{
			// token is not changed, so we can return it as it was passed.
			Token:     strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer "),
			ID:        sess.ID,
			AccountID: sess.AccountID,
			CreatedAt: sess.CreatedAt,
			ExpiresAt: sess.ExpiresAt,
		},
		nil,
	))
}

type getSessionSchema struct {
	Token     string    `json:"token"`
	ID        string    `json:"id"`
	AccountID int64     `json:"accountID"`
	CreatedAt time.Time `json:"createdAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}
