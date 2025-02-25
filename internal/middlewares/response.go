package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

// SetContentLengthMiddleware устанавливает заголовок Content-Length.
func ContentLengthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rr := &responseRecorder{ResponseWriter: w}
		next.ServeHTTP(rr, r)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", rr.contentLength))
	})
}

// SetDateMiddleware устанавливает заголовок Date в формате RFC1123.
func DateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Date", time.Now().Format(time.RFC1123))
		next.ServeHTTP(w, r)
	})
}

// SetTextPlainContentType устанавливает заголовок Content-Type для текстовых данных.
func TextPlainContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// SetJSONContentType устанавливает заголовок Content-Type для JSON.
func JSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// SetHTMLContentType устанавливает заголовок Content-Type для HTML.
func HTMLContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

// responseRecorder записывает статус код и длину содержимого для последующего использования.
type responseRecorder struct {
	http.ResponseWriter
	statusCode    int
	contentLength int
}

// Write переопределяет метод Write для записи содержимого и подсчета длины.
func (rr *responseRecorder) Write(b []byte) (int, error) {
	rr.contentLength += len(b)
	return rr.ResponseWriter.Write(b)
}
