package middleware

import (
	"log"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// LoggingMiddleware logs method, path, status and duration for each request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		recorder := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		start := time.Now()

		next.ServeHTTP(recorder, req)

		log.Printf("%s %s -> %d (%s)", req.Method, req.URL.Path, recorder.status, time.Since(start))
	})
}
