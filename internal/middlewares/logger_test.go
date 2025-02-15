package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Define a mock logger struct to simulate the logger's behavior
type MockLogger struct {
	InfoCalled bool
	InfoArgs   string
}

// Implement the Info method for the MockLogger
func (m *MockLogger) Info(msg string, fields ...zap.Field) {
	m.InfoCalled = true
	m.InfoArgs = msg
}

// Implement the Debug method for the MockLogger (not used in this test, but defined for completeness)
func (m *MockLogger) Debug(msg string, fields ...zap.Field) {}

// Implement the Error method for the MockLogger (not used in this test, but defined for completeness)
func (m *MockLogger) Error(msg string, fields ...zap.Field) {}

// TestLoggerMiddlewareInfoCall tests that the Info method of the logger is called correctly.
func TestLoggerMiddlewareInfoCall(t *testing.T) {
	// Define the mock logger instance
	mockLogger := &MockLogger{
		InfoCalled: false,
		InfoArgs:   "",
	}

	// Create the Gin engine
	r := gin.Default()

	// Use the LoggerMiddleware with the mock logger
	r.Use(LoggerMiddleware(mockLogger))

	// Define a test route
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	// Create a test request (using gin.CreateTestContext to simulate requests)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	// Simulate the request
	r.ServeHTTP(w, req)

	// Simulate request processing time
	time.Sleep(50 * time.Millisecond)

	// Assert that the Info method was called
	assert.True(t, mockLogger.InfoCalled, "Info method should have been called")

	// Assert that the Info message contains "Request processed"
	assert.Contains(t, mockLogger.InfoArgs, "Request processed", "Info should log 'Request processed'")
}
