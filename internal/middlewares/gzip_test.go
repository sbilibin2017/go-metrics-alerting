package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Простая обработка для тестирования middleware.
func simpleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

// Функция для сжатия данных в GZIP.
func gzipCompress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, err
	}
	gzipWriter.Close()
	return buf.Bytes(), nil
}

func TestGzipMiddleware_CompressResponse(t *testing.T) {
	// Создаем новый HTTP запрос
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	// Создаем тестовый респондер и проксируем запрос через middleware
	rr := httptest.NewRecorder()
	handler := GzipMiddleware(http.HandlerFunc(simpleHandler))

	// Имитируем, что клиент поддерживает gzip
	req.Header.Set("Accept-Encoding", "gzip")

	// Обрабатываем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем, что ответ имеет заголовок Content-Encoding: gzip
	assert.Equal(t, "gzip", rr.Header().Get("Content-Encoding"))

	// Проверяем, что тело ответа сжато (попробуем распаковать)
	respBody := rr.Body.Bytes()
	gzipReader, err := gzip.NewReader(bytes.NewReader(respBody))
	assert.NoError(t, err)

	decompressedData := &bytes.Buffer{}
	_, err = decompressedData.ReadFrom(gzipReader)
	assert.NoError(t, err)

	// Проверяем, что распакованные данные соответствуют ожидаемым
	assert.Equal(t, "Hello, World!", decompressedData.String())
}

func TestGzipMiddleware_DecompressRequest(t *testing.T) {
	// Создаем сжатыми данные для запроса
	originalData := []byte("Hello, GZIP request!")
	compressedData, err := gzipCompress(originalData)
	assert.NoError(t, err)

	// Создаем новый HTTP запрос с сжатыми данными в теле
	req, err := http.NewRequest("POST", "/", bytes.NewReader(compressedData))
	assert.NoError(t, err)

	// Устанавливаем заголовок Content-Encoding: gzip
	req.Header.Set("Content-Encoding", "gzip")

	// Создаем тестовый респондер
	rr := httptest.NewRecorder()
	handler := GzipMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверяем, что тело запроса распаковано
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		assert.Equal(t, originalData, body)
		w.Write([]byte("Request received"))
	}))

	// Обрабатываем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем, что ответ пришел правильно
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Request received", rr.Body.String())
}
