package logger

import (
	"go.uber.org/zap"
)

// Глобальная переменная для хранения логгера
var Logger *zap.Logger

// init — инициализация глобального логгера с уровнем DEBUG.
func init() {
	// Создаем конфигурацию для разработки (вы можете использовать другую конфигурацию по мере необходимости)
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	// Строим логгер
	var err error
	Logger, err = zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic("Failed to initialize zap logger: " + err.Error())
	}

	// Логируем информацию о том, что логгер инициализирован
	Logger.Debug("Logger is initialized with level DEBUG")
}
