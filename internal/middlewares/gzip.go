package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware — это middleware, которое сжимает ответы в формате GZIP и распаковывает запросы с GZIP.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gzipReader, _ := gzip.NewReader(r.Body)
			defer gzipReader.Close()
			r.Body = io.NopCloser(gzipReader)
		}
		grw := &gzipResponseWriter{ResponseWriter: w}
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			grw.gzipEnabled = true
			w.Header().Set("Content-Encoding", "gzip")
		}
		next.ServeHTTP(grw, r)
	})
}

// gzipResponseWriter — это обёртка для ResponseWriter для сжатия ответа.
type gzipResponseWriter struct {
	http.ResponseWriter
	gzipEnabled bool
}

// Write записывает данные в ответ. Если сжатие включено, то сжимаем данные перед отправкой.
func (grw *gzipResponseWriter) Write(p []byte) (int, error) {
	if grw.gzipEnabled {
		gzipWriter := gzip.NewWriter(grw.ResponseWriter)
		defer gzipWriter.Close()
		return gzipWriter.Write(p)
	}
	return grw.ResponseWriter.Write(p)
}
