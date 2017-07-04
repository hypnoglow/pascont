package middleware

import (
	"log"
	"net/http"
)

func Recover(next http.Handler, errorLogger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errorLogger.Printf("panic: %+v", err)
				w.WriteHeader(500)
			}
		}()

		next.ServeHTTP(w, req)
	})
}
