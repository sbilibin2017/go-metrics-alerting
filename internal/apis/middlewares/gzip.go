package middlewares

import (
	"compress/gzip"
	"go-metrics-alerting/internal/apis/responses"
	"net/http"
	"strings"
)

// GzipMiddleware compresses the response body if the client accepts gzip encoding
// and decompresses the request body if the client sends gzip-encoded data.
func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle Gzip Decompression on Request Body
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			// Decompress the body if it's gzip compressed
			gzipReader, err := gzip.NewReader(r.Body)
			if err != nil {
				responses.InternalServerErrorResponse(w, err)
				return
			}
			defer gzipReader.Close()
			// Set the new request body as the decompressed reader
			r.Body = gzipReader
		}

		// Create a response writer wrapper that will handle Gzip compression for the response
		gzipResponseWriter := &GzipResponseWriter{ResponseWriter: w}

		// Check if the client supports Gzip
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// If the client accepts Gzip encoding, wrap the ResponseWriter
			w.Header().Set("Content-Encoding", "gzip")
			gzipResponseWriter.Compress = true
		}

		// Pass the request to the next handler
		next.ServeHTTP(gzipResponseWriter, r)
	})
}

// GzipResponseWriter wraps the http.ResponseWriter to handle Gzip compression
type GzipResponseWriter struct {
	http.ResponseWriter
	Compress bool
	Writer   *gzip.Writer
}

func (w *GzipResponseWriter) Write(p []byte) (n int, err error) {
	if w.Compress {
		// If compression is enabled, write to the gzip writer
		if w.Writer == nil {
			w.Writer = gzip.NewWriter(w.ResponseWriter)
		}
		defer w.Writer.Close()
		return w.Writer.Write(p)
	}
	// Otherwise, write the data normally
	return w.ResponseWriter.Write(p)
}

func (w *GzipResponseWriter) WriteHeader(statusCode int) {
	if w.Compress {
		// If compression is enabled, write headers with gzip encoding
		w.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	w.ResponseWriter.WriteHeader(statusCode)
}
