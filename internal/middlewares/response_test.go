package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContentLengthMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// создаем новый сервер с middleware
	tt := httptest.NewServer(ContentLengthMiddleware(handler))
	defer tt.Close()

	res, err := http.Get(tt.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "13", res.Header.Get("Content-Length"))
}

func TestDateMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// создаем новый сервер с middleware
	tt := httptest.NewServer(DateMiddleware(handler))
	defer tt.Close()

	res, err := http.Get(tt.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	dateHeader := res.Header.Get("Date")
	// Проверяем, что дата в заголовке существует и соответствует формату RFC1123
	_, err = time.Parse(time.RFC1123, dateHeader)
	assert.NoError(t, err)
}

func TestTextPlainContentTypeMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	// создаем новый сервер с middleware
	tt := httptest.NewServer(TextPlainContentType(handler))
	defer tt.Close()

	res, err := http.Get(tt.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
}

func TestJSONContentTypeMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"message": "Hello, world!"}`))
	})

	// создаем новый сервер с middleware
	tt := httptest.NewServer(JSONContentType(handler))
	defer tt.Close()

	res, err := http.Get(tt.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", res.Header.Get("Content-Type"))
}

func TestHTMLContentTypeMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<html><body>Hello, world!</body></html>"))
	})

	// создаем новый сервер с middleware
	tt := httptest.NewServer(HTMLContentType(handler))
	defer tt.Close()

	res, err := http.Get(tt.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "text/html; charset=utf-8", res.Header.Get("Content-Type"))
}
