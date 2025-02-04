package apiclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Мок-сервер для эмуляции ответов API
func newMockServer(responseCode int, responseBody string, delay time.Duration) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(responseCode)
		w.Write([]byte(responseBody))
	})
	return httptest.NewServer(handler)
}

// Тест: успешный GET-запрос
func TestApiClientEngine_Get_Success(t *testing.T) {
	mockServer := newMockServer(http.StatusOK, `{"message": "ok"}`, 0)
	defer mockServer.Close()

	client := NewClient(ApiClientOptions{Timeout: 2 * time.Second, RetryCount: 1})
	resp, err := client.Get(mockServer.URL, nil, nil)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.JSONEq(t, `{"message": "ok"}`, string(resp.Body))
}

// Тест: GET-запрос с тайм-аутом
func TestApiClientEngine_Get_Timeout(t *testing.T) {
	mockServer := newMockServer(http.StatusOK, `{"message": "ok"}`, 3*time.Second)
	defer mockServer.Close()

	client := NewClient(ApiClientOptions{Timeout: 1 * time.Second, RetryCount: 1})
	resp, err := client.Get(mockServer.URL, nil, nil)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrTimeout)
}

// Тест: GET-запрос с некорректным ответом
func TestApiClientEngine_Get_InvalidResponse(t *testing.T) {
	mockServer := newMockServer(http.StatusInternalServerError, ``, 0)
	defer mockServer.Close()

	client := NewClient(ApiClientOptions{Timeout: 2 * time.Second, RetryCount: 1})
	resp, err := client.Get(mockServer.URL, nil, nil)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrInvalidResponse)
}

// Тест: успешный POST-запрос
func TestApiClientEngine_Post_Success(t *testing.T) {
	mockServer := newMockServer(http.StatusCreated, `{"message": "created"}`, 0)
	defer mockServer.Close()

	client := NewClient(ApiClientOptions{Timeout: 2 * time.Second, RetryCount: 1})
	resp, err := client.Post(mockServer.URL, nil, `{"data": "test"}`, nil)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.JSONEq(t, `{"message": "created"}`, string(resp.Body))
}

// Тест: POST-запрос с тайм-аутом
func TestApiClientEngine_Post_Timeout(t *testing.T) {
	mockServer := newMockServer(http.StatusOK, `{"message": "ok"}`, 3*time.Second)
	defer mockServer.Close()

	client := NewClient(ApiClientOptions{Timeout: 1 * time.Second, RetryCount: 1})
	resp, err := client.Post(mockServer.URL, nil, `{"data": "test"}`, nil)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrTimeout)
}

// Тест: POST-запрос с некорректным ответом
func TestApiClientEngine_Post_InvalidResponse(t *testing.T) {
	mockServer := newMockServer(http.StatusInternalServerError, ``, 0)
	defer mockServer.Close()

	client := NewClient(ApiClientOptions{Timeout: 2 * time.Second, RetryCount: 1})
	resp, err := client.Post(mockServer.URL, nil, `{"data": "test"}`, nil)

	assert.Nil(t, resp)
	assert.ErrorIs(t, err, ErrInvalidResponse)
}

// Тест: обработка ошибки тайм-аута
func TestIsTimeoutError(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	<-ctx.Done()
	err := ctx.Err()

	assert.True(t, isTimeoutError(err))
	assert.False(t, isTimeoutError(errors.New("random error")))
}
