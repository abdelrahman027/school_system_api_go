package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start)
		wrappedWriter.Header().Set("X-Response-Time", duration.String())
		// Log the response time
		fmt.Printf("Method : %v , url: %v ,code: %v ,duration: %v\n", r.Method, r.URL.Path, wrappedWriter.statusCode, duration.String())

		// w.Header().Set("X-Response-Time", duration.String())
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
