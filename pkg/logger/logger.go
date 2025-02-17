package logger

import (
	"go.uber.org/zap"
)

// Глобальный zap.Logger
var Log *zap.Logger

// init — инициализация глобального логгера с уровнем DEBUG.
func init() {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	Log, _ = zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	Log.Debug("Logger is initialized with level DEBUG")
}
