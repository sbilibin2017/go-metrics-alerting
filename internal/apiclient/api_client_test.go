package apiclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPost_AddSchemeIfMissing тестирует добавление схемы "http://" по умолчанию, если в URL нет схемы
func TestPost_AddSchemeIfMissing(t *testing.T) {
	// Создаем тестовый сервер, который будет отвечать на POST-запросы
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Проверим, что запрос был POST
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		// Отправим успешный ответ
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "success"}`)
	})

	// Создаем сервер на базе этого обработчика
	server := httptest.NewServer(handler)
	defer server.Close()

	// Создаем экземпляр клиента Resty
	client := NewRestyClient()

	// Тестируем URL без схемы
	urlWithoutScheme := server.URL[len("http://"):] // Просто берем URL без "http://"
	urlWithoutScheme = "http://" + urlWithoutScheme // Добавляем схему

	// Выполняем запрос с правильным URL
	resp, err := client.Post(urlWithoutScheme+"/test", map[string]string{})

	// Проверка, что URL был исправлен и запрос отправлен правильно
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `{"message": "success"}`, resp.Body)
}

// TestPost_SetHeaders проверяет установку заголовков в запрос
func TestPost_SetHeaders(t *testing.T) {
	// Создаем тестовый сервер, который будет отвечать на POST-запросы
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Проверим заголовки
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))

		// Отправим успешный ответ
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "success"}`)
	})

	// Создаем сервер на базе этого обработчика
	server := httptest.NewServer(handler)
	defer server.Close()

	// Создаем экземпляр клиента Resty
	client := NewRestyClient()

	// Заголовки для запроса
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}

	// Выполняем запрос с установленными заголовками
	resp, err := client.Post(server.URL+"/test", headers)

	// Проверка результатов
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `{"message": "success"}`, resp.Body)
}

// TestPost_HeadersWithNoContentType проверяет, что заголовок Content-Type установлен по умолчанию, если его нет в запросе
func TestPost_HeadersWithNoContentType(t *testing.T) {
	// Создаем тестовый сервер, который будет отвечать на POST-запросы
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Проверим, что Content-Type был установлен на default "application/json"
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// Отправим успешный ответ
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "success"}`)
	})

	// Создаем сервер на базе этого обработчика
	server := httptest.NewServer(handler)
	defer server.Close()

	// Создаем экземпляр клиента Resty
	client := NewRestyClient()

	// Отправляем запрос без Content-Type
	headers := map[string]string{
		"Authorization": "Bearer token",
	}

	// Выполняем запрос с заголовками
	resp, err := client.Post(server.URL+"/test", headers)

	// Проверка результатов
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `{"message": "success"}`, resp.Body)
}

// TestPost_Success тестирует успешный POST-запрос
func TestPost_Success(t *testing.T) {
	// Создаем тестовый сервер, который будет отвечать на POST-запросы
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Проверим, что запрос был POST
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		// Отправим успешный ответ
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, `{"message": "success"}`)
	})

	// Создаем сервер на базе этого обработчика
	server := httptest.NewServer(handler)
	defer server.Close()

	// Создаем экземпляр клиента Resty
	client := NewRestyClient()

	// Выполняем тестируемую функцию
	resp, err := client.Post(server.URL+"/test", map[string]string{})

	// Проверка результатов
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, `{"message": "success"}`, resp.Body)
}

// TestPost_Failure тестирует ситуацию с ошибкой при выполнении POST-запроса
func TestPost_Failure(t *testing.T) {
	// Создаем тестовый сервер, который будет генерировать ошибку
	handler := http.NewServeMux()
	handler.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Симулируем ошибку
		http.Error(w, "Bad Request", http.StatusBadRequest)
	})

	// Создаем сервер на базе этого обработчика
	server := httptest.NewServer(handler)
	defer server.Close()

	// Создаем экземпляр клиента Resty
	client := NewRestyClient()

	// Выполняем тестируемую функцию
	resp, err := client.Post(server.URL+"/test", map[string]string{})

	// Проверка результатов
	assert.Error(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 400, resp.StatusCode)
	assert.Equal(t, "Bad Request", resp.Body)
}

// TestPost_InvalidURL тестирует ситуацию с некорректным URL
func TestPost_InvalidURL(t *testing.T) {
	// Создаем сервер, который не будет запущен
	invalidURL := "http://localhost:9999/nonexistent"

	// Создаем экземпляр клиента Resty
	client := NewRestyClient()

	// Выполняем тестируемую функцию
	resp, err := client.Post(invalidURL, map[string]string{})

	// Проверка результатов
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "error while making POST request")
}
