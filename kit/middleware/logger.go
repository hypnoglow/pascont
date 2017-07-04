package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseRecorder struct {
	http.ResponseWriter
	status        int
	statusWritten bool
}

func newResponseRecorder(origin http.ResponseWriter) *responseRecorder {
	return &responseRecorder{origin, http.StatusOK, false}
}

func (rr *responseRecorder) WriteHeader(status int) {
	if !rr.statusWritten {
		rr.status = status
	}
	rr.ResponseWriter.WriteHeader(status)
}

func (rr responseRecorder) Status() int {
	return rr.status
}

func Logger(next http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rr := newResponseRecorder(w)

		before := time.Now()

		next.ServeHTTP(rr, req)

		//dur := time.Since(before).Nanoseconds() / 1e6
		logger.Printf("%s %s : %d %d ms", req.Method, req.URL.String(), rr.Status(), time.Since(before))
	})
}
