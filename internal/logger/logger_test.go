package logger

import (
	"testing"

	"go.uber.org/zap/zaptest"
)

func TestLoggerInitialization(t *testing.T) {
	// Проверяем, что глобальная переменная Log не равна nil после инициализации
	if Log == nil {
		t.Fatal("Logger is not initialized")
	}
}

func TestLoggerLevel(t *testing.T) {
	// Используем zaptest, чтобы проверить уровень логирования
	// zaptest.NewLogger создает новый логгер, который можно тестировать
	logger := zaptest.NewLogger(t)

	// Логируем сообщение с уровнем DEBUG
	logger.Debug("Testing debug level")

	// Здесь можно добавить логику для проверки выводов или поведения логирования.
	// В данном случае мы проверим, что ничего не сломается.
}
