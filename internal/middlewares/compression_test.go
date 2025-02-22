package middlewares

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompressionMiddleware(t *testing.T) {
	// Тестовые данные для сжатия
	data := []byte("Hello, GZIP Compression!")

	// Создаем запрос
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	// Создаем ответ
	rr := httptest.NewRecorder()

	// Обработчик с Middleware для сжатия
	handler := CompressionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Записываем тестовые данные в ответ
		w.Write(data)
	}))

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Получаем результат
	resp := rr.Result()
	defer resp.Body.Close()

	// Проверяем, что ответ сжат с использованием GZIP
	assert.Equal(t, "gzip", resp.Header.Get("Content-Encoding"))

	// Читаем сжатые данные из ответа
	compressedBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Декодируем сжатые данные
	gzipReader, err := gzip.NewReader(bytes.NewReader(compressedBody))
	require.NoError(t, err)
	defer gzipReader.Close()

	// Распаковываем данные
	decompressedData, err := io.ReadAll(gzipReader)
	require.NoError(t, err)

	// Проверяем, что распакованные данные совпадают с исходными
	assert.Equal(t, data, decompressedData)
}

func TestDecompressionMiddleware(t *testing.T) {
	// Тестовые данные для сжатия
	data := []byte("Hello, GZIP Decompression!")

	// Сжимаем тестовые данные для имитации запроса с GZIP-данными
	compressedData, err := compress(data)
	require.NoError(t, err)

	// Создаем запрос с сжатыми данными в теле
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(compressedData))
	req.Header.Set("Content-Encoding", "gzip")
	// Создаем ответ
	rr := httptest.NewRecorder()

	// Обработчик с Middleware для разжатия
	handler := DecompressionMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Читаем тело запроса после разжатия
		decompressedBody, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		// Проверяем, что разжатое тело совпадает с исходными данными
		assert.Equal(t, data, decompressedBody)
	}))

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем, что статус ответа корректен
	assert.Equal(t, http.StatusOK, rr.Code)
}
