package middleware

import (
	"context"
	"net/http"
	"strings"
)

type ContextKeyPathID struct{}

func PathID(next http.Handler, pathPrefix string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req = req.WithContext(context.WithValue(
			req.Context(),
			ContextKeyPathID{},
			strings.TrimPrefix(req.URL.Path, pathPrefix),
		))

		next.ServeHTTP(w, req)
	})
}
