package middlewares

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

// Compress сжимает данные с помощью GZIP.
func compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gzipWriter := gzip.NewWriter(&b)
	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed to compress data: %v", err)
	}
	err = gzipWriter.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close gzip writer: %v", err)
	}
	return b.Bytes(), nil
}

// Decompress распаковывает данные, сжатые с помощью GZIP.
func decompress(data []byte) ([]byte, error) {
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzipReader.Close()
	var b bytes.Buffer
	_, err = b.ReadFrom(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to decompress data: %v", err)
	}

	return b.Bytes(), nil
}

// DecompressRequestBody разжимает тело запроса
func decompressRequestBody(body io.Reader) ([]byte, error) {
	compressedData, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %v", err)
	}
	return decompress(compressedData)
}

// CompressionMiddleware - middleware для сжатия данных в ответе с использованием GZIP.
func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		originalWriter := w
		gzipWriter := &gzipResponseWriter{ResponseWriter: originalWriter}
		next.ServeHTTP(gzipWriter, r)
		compressedData, err := compress(gzipWriter.body.Bytes())
		if err != nil {
			http.Error(w, "failed to compress response data", http.StatusInternalServerError)
			return
		}
		originalWriter.Header().Set("Content-Encoding", "gzip")
		originalWriter.Header().Set("Content-Length", fmt.Sprintf("%d", len(compressedData)))
		originalWriter.Write(compressedData)
	})
}

// gzipResponseWriter - обертка для ResponseWriter, чтобы захватить данные ответа
type gzipResponseWriter struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (gw *gzipResponseWriter) Write(p []byte) (int, error) {
	return gw.body.Write(p)
}

// DecompressionMiddleware - middleware для разжатия данных в запросах с использованием GZIP.
func DecompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			decompressedBody, err := decompressRequestBody(r.Body)
			if err != nil {
				http.Error(w, "failed to decompress request data", http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(decompressedBody))
		}
		next.ServeHTTP(w, r)
	})
}
