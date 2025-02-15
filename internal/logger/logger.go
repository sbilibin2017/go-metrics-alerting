package logger

import (
	"os"

	"go.uber.org/zap"
)

// Переменная для хранения логгера
var logger *zap.Logger

// init конфигурирует логгер на основе переменных окружения.
func init() {
	// Получаем уровень логирования из переменной окружения
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "INFO" // По умолчанию INFO
	}

	// Конфигурируем логгер для JSON-формата
	zapConfig := zap.NewProductionConfig()

	// Устанавливаем уровень логирования в зависимости от конфигурации
	switch logLevel {
	case "DEBUG":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "ERROR":
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Создаем логгер с указанной конфигурацией
	logger, _ = zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1)) // Добавляем информацию о месте вызова

}

// Debug логирует сообщения уровня Debug.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info логирует сообщения уровня Info.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Error логирует сообщения уровня Error.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
