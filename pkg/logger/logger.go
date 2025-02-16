package logger

import (
	"go.uber.org/zap"
)

// Logger is an implementation of the Logger interface using zap.
type Logger struct {
	logger *zap.Logger
}

// NewLogger creates a new Logger instance based on the provided configuration.
func NewLogger(config *LoggerConfig) (*Logger, error) {
	var zapConfig zap.Config

	// Set the zap config based on the log level from the config
	switch config.LogLevel {
	case DEBUG:
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case ERROR:
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default: // INFO as default
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	// Create the zap logger with the configured settings
	logger, err := zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return &Logger{logger: logger}, nil
}

// Info logs a message at the Info level.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Debug logs a message at the Debug level.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Error logs a message at the Error level.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}
