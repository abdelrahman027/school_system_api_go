package middlewares

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the client supports gzip encoding
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}
		// Set the appropriate headers
		w.Header().Set("Content-Encoding", "gzip")

		// Create a gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()
		w = &gzipResponseWriter{ResponseWriter: w, writer: gz}

		// Call the next handler with the gzip writer

		next.ServeHTTP(w, r)
	})
}

type gzipResponseWriter struct {
	http.ResponseWriter
	writer *gzip.Writer
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}
