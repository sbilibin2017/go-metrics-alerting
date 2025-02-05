package logger

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLoggerInitialization(t *testing.T) {
	// Проверяем, что глобальный Logger инициализирован
	assert.NotNil(t, Logger, "Logger should be initialized")

	// Проверяем, что у Logger установлен уровень DebugLevel
	assert.Equal(t, logrus.DebugLevel, Logger.Level, "Logger should be set to Debug level")

	// Проверяем, что у Logger установлен JSONFormatter
	_, ok := Logger.Formatter.(*logrus.JSONFormatter)
	assert.True(t, ok, "Logger should use JSONFormatter")
}

func TestLoggerOutput(t *testing.T) {
	// Захватываем вывод логгера в буфер
	var buf bytes.Buffer
	Logger.SetOutput(&buf)

	// Логируем сообщение
	Logger.Info("Test message")

	// Проверяем, что в буфере есть строка "Test message"
	logOutput := buf.String()

	// Проверяем, что в выводе содержится сообщение
	assert.Contains(t, logOutput, `"level":"info"`)
	assert.Contains(t, logOutput, `"msg":"Test message"`)
}

func TestLoggerLevel(t *testing.T) {
	// Создаем новый логгер с уровнем логирования "info"
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)

	// Захватываем вывод
	var buf bytes.Buffer
	logger.SetOutput(&buf)

	// Логируем сообщения на разных уровнях
	logger.Debug("This is a debug message") // Это сообщение не должно быть выведено
	logger.Info("This is an info message")  // Это сообщение должно быть выведено

	// Проверяем, что в выводе есть только сообщение уровня Info
	logOutput := buf.String()
	assert.NotContains(t, logOutput, "This is a debug message")
	assert.Contains(t, logOutput, "This is an info message")
}
