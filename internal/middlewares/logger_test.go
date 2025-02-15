package middlewares

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLoggerMiddleware(t *testing.T) {
	// Создаем новый логгер для тестов
	logger, _ := zap.NewDevelopment()

	// Создаем новый роутер и подключаем middleware
	r := gin.New()
	r.Use(LoggerMiddleware(logger))

	// Определяем тестовый эндпоинт
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	// Создаем новый запрос
	req, _ := http.NewRequest("GET", "/test", nil)

	// Создаем рекордер для захвата ответа
	rr := httptest.NewRecorder()

	// Отправляем запрос
	r.ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	assert.Equal(t, "test", rr.Body.String())
}

func TestLoggerMiddleware_LogRequest(t *testing.T) {
	// Создаем буфер для захвата логов
	var logOutput bytes.Buffer

	// Создаем конфигурацию для логгера с буфером
	writeSyncer := zapcore.AddSync(&logOutput)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	// Создаем логгер с конфигурацией
	logger := zap.New(core)
	defer logger.Sync() // Отсроченное сбрасывание

	// Создаем новый роутер и подключаем middleware
	r := gin.New()
	r.Use(LoggerMiddleware(logger))

	// Определяем тестовый эндпоинт
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	// Создаем новый запрос
	req, _ := http.NewRequest("GET", "/test", nil)

	// Создаем рекордер для захвата ответа
	rr := httptest.NewRecorder()

	// Отправляем запрос
	r.ServeHTTP(rr, req)

	// Проверяем, что лог был записан
	logContent := logOutput.String()
	assert.Contains(t, logContent, `"method":"GET"`)
	assert.Contains(t, logContent, `"uri":"/test"`) // Теперь правильно ожидается
	assert.Contains(t, logContent, `"status_code":200`)
	assert.Contains(t, logContent, `"content_length":4`) // "test" имеет длину 4 символа
	assert.Contains(t, logContent, `"duration"`)
}

func TestLoggerMiddleware_LogRequestDetails(t *testing.T) {
	// Создаем буфер для захвата логов
	var logOutput bytes.Buffer

	// Создаем конфигурацию для логгера с буфером
	writeSyncer := zapcore.AddSync(&logOutput)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)

	// Создаем логгер с конфигурацией
	logger := zap.New(core)
	defer logger.Sync() // Отсроченное сбрасывание

	// Создаем новый роутер и подключаем middleware
	r := gin.New()
	r.Use(LoggerMiddleware(logger))

	// Определяем тестовый эндпоинт
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	// Создаем новый запрос
	req, _ := http.NewRequest("GET", "/test", nil)

	// Создаем рекордер для захвата ответа
	rr := httptest.NewRecorder()

	// Отправляем запрос
	r.ServeHTTP(rr, req)

	// Проверяем, что лог был записан с корректными значениями
	logContent := logOutput.String()
	assert.Contains(t, logContent, `"method":"GET"`)
	assert.Contains(t, logContent, `"uri":"/test"`)
	assert.Contains(t, logContent, `"status_code":200`)
	assert.Contains(t, logContent, `"content_length":4`) // "test" имеет длину 4 символа
	assert.Contains(t, logContent, `"duration"`)
}
