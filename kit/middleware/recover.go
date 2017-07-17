package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

func Recover(next http.Handler, errorLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorLogger.Printf("panic: %v %s", err, debug.Stack())
				w.WriteHeader(500)
			}
		}()

		next.ServeHTTP(w, req)
	})
}
