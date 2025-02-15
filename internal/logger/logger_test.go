package logger

import (
	"testing"

	"go-metrics-alerting/internal/configs"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockLogger is used to mock the logger methods
type MockLogger struct {
	mock.Mock
}

// Info mocks the Info method
func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

// Debug mocks the Debug method
func (m *MockLogger) Debug(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

// Error mocks the Error method
func (m *MockLogger) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func TestLoggerWithInfoLevel(t *testing.T) {

	// Create LoggerConfig with INFO level
	config := &configs.LoggerConfig{
		LogLevel: configs.INFO,
	}

	// Create the logger with the configuration
	log, err := NewLogger(config)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Call the methods on the logger
	log.Info("This is an info message")
	log.Debug("This is a debug message")
	log.Error("This is an error message")

}
