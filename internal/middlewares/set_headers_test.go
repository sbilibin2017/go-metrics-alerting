package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSetHeadersMiddleware(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		expectedContent string
		expectedType    string
		expectedDate    string
		expectedLength  int
	}{
		{
			name:            "Test JSON content",
			url:             "/",
			expectedContent: `{"message": "Hello, World!"}`,
			expectedType:    "application/json",
			expectedDate:    time.Now().UTC().Format(time.RFC1123),
			expectedLength:  len(`{"message": "Hello, World!"}`),
		},
		{
			name:            "Test Plain Text content",
			url:             "/text",
			expectedContent: "Hello, World in plain text!",
			expectedType:    "text/plain",
			expectedDate:    time.Now().UTC().Format(time.RFC1123),
			expectedLength:  len("Hello, World in plain text!"),
		},
		{
			name:            "Test HTML content",
			url:             "/html",
			expectedContent: "<html><body><h1>Hello, World in HTML!</h1></body></html>",
			expectedType:    "text/html",
			expectedDate:    time.Now().UTC().Format(time.RFC1123),
			expectedLength:  len("<html><body><h1>Hello, World in HTML!</h1></body></html>"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			// Создаем ResponseRecorder, чтобы захватить ответ
			rr := httptest.NewRecorder()

			// Обработчик, который возвращает разные ответы в зависимости от URL
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/text" {
					w.Write([]byte("Hello, World in plain text!"))
				} else if r.URL.Path == "/html" {
					w.Write([]byte("<html><body><h1>Hello, World in HTML!</h1></body></html>"))
				} else {
					w.Write([]byte(`{"message": "Hello, World!"}`))
				}
			})

			// Применяем наш middleware
			middleware := SetHeadersMiddleware(handler)
			middleware.ServeHTTP(rr, req)

			// Проверяем заголовки
			assert.Equal(t, tt.expectedType, rr.Header().Get("Content-Type"))
			assert.Contains(t, rr.Header().Get("Date"), tt.expectedDate[:len(tt.expectedDate)-1]) // Проверяем дату с точностью до секунды
			assert.Equal(t, fmt.Sprintf("%d", tt.expectedLength), rr.Header().Get("Content-Length"))

			// Проверяем контент
			assert.Equal(t, tt.expectedContent, rr.Body.String())
		})
	}
}
