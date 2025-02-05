package apiclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRestyClient_Post(t *testing.T) {
	// Создаем тестовый сервер
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))
	defer ts.Close()

	// Создаем клиента
	client := NewRestyClient()

	// Отправляем POST-запрос с URL, который уже имеет схему (http://)
	resp, err := client.Post(ts.URL, map[string]string{"Content-Type": "application/json"})

	// Проверяем, что ошибок нет
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Success", resp.Body)
}

func TestRestyClient_Post_WithoutScheme(t *testing.T) {
	// Создаем тестовый сервер, который обрабатывает путь /update
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка метода и пути
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
			return
		}
		if r.URL.Path != "/update" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	}))
	defer ts.Close()

	// Создаем клиента
	client := NewRestyClient()

	// Отправляем POST-запрос на путь /update
	resp, err := client.Post(ts.URL+"/update", map[string]string{"Content-Type": "application/json"})

	// Проверяем, что ошибок нет
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "Success", resp.Body)
}

func TestRestyClient_Post_ServerError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer ts.Close()

	client := NewRestyClient()
	resp, err := client.Post(ts.URL, nil)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, "Internal Server Error", resp.Body)
}

func TestRestyClient_Post_InvalidURL(t *testing.T) {
	client := NewRestyClient()
	resp, err := client.Post("http://invalid-url", nil)

	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestAddSchemeIfMissing(t *testing.T) {
	tests := []struct {
		urlString     string
		expectedURL   string
		expectedError bool
	}{
		{
			urlString:     "localhost:8080",        // Нет схемы, должна быть добавлена
			expectedURL:   "http://localhost:8080", // Добавлено http://
			expectedError: false,
		},
		{
			urlString:     "http://localhost:8080", // Уже есть схема http://
			expectedURL:   "http://localhost:8080", // Схема не изменится
			expectedError: false,
		},
		{
			urlString:     "https://example.com", // Уже есть схема https://
			expectedURL:   "https://example.com", // Схема не изменится
			expectedError: false,
		},
		{
			urlString:     "invalid-url", // Некорректный URL
			expectedURL:   "",            // Ошибка
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.urlString, func(t *testing.T) {
			result, err := addSchemeIfMissing(tt.urlString)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, "", result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedURL, result)
			}
		})
	}
}
