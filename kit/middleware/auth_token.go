package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pborman/uuid"

	"github.com/hypnoglow/pascont/kit"
	"github.com/hypnoglow/pascont/kit/schema"
)

type ContextKeySessionID struct{}

func AuthToken(next http.Handler, extractor func(token string) (id string, err error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		token := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			next.ServeHTTP(w, req)
			return
		}

		sessionID, err := extractor(token)
		if err != nil || uuid.Parse(sessionID) == nil {
			kit.RespondWithError(w, 401, schema.ErrorFromMessage(
				"Authorization token provided is invalid.",
			))
			return
		}

		req = req.WithContext(context.WithValue(
			req.Context(),
			ContextKeySessionID{},
			sessionID,
		))

		next.ServeHTTP(w, req)
	})
}
