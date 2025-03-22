package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

// GzipMiddleware logs incoming requests and outgoing responses, compressing or decompressing the body as necessary.
func GzipMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// If the request is gzip encoded, decompress the body
			if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {

				// Decompress request body
				reader, err := decompress(r.Body)
				if err != nil {

					http.Error(w, "Error reading gzip body", http.StatusBadRequest)
					return
				}

				r.Body = io.NopCloser(reader)
			}

			// If the client accepts gzip, prepare to compress the response
			if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {

				// Set the Content-Encoding header to indicate the response is gzipped
				w.Header().Set("Content-Encoding", "gzip")

				// Attempt to compress the response body
				// Create a custom writer to capture the compressed response
				gzipWriter := gzip.NewWriter(w)
				defer gzipWriter.Close()

				// Capture the original response writer
				originalWriter := w
				w = &responsezWriter{
					ResponseWriter: originalWriter,
					Writer:         gzipWriter,
				}

				// Proceed to the next handler and write the compressed response
				next.ServeHTTP(w, r)
				return
			}

			// Proceed to the next handler if no compression is needed
			next.ServeHTTP(w, r)
		})
	}
}

// Compress compresses the data from the `data` reader and writes the compressed output to the `writer`.
func compress(writer io.Writer, data io.Reader) error {
	// Create a new gzip writer
	gzipWriter := gzip.NewWriter(writer)
	defer gzipWriter.Close()

	// Copy the data into the gzip writer
	_, err := io.Copy(gzipWriter, data)
	return err
}

// Decompress decompresses the data from the `reader` and returns the decompressed data as a new reader.
func decompress(reader io.Reader) (io.Reader, error) {
	// Create a new gzip reader
	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}

	return gzipReader, nil
}

// Custom responseWriter to capture compressed responses
type responsezWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (rw *responsezWriter) WriteHeader(statusCode int) {
	if rw.Writer == nil {
		// Prevent calling WriteHeader multiple times
		rw.ResponseWriter.WriteHeader(statusCode)
		return
	}
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responsezWriter) Write(p []byte) (int, error) {
	if rw.Writer != nil {
		return rw.Writer.Write(p)
	}
	return rw.ResponseWriter.Write(p)
}
