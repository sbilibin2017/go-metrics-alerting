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

// TestLoggerMiddleware checks that the LoggerMiddleware works as expected
func TestLoggerMiddleware(t *testing.T) {
	// Create a new Gin router and add LoggerMiddleware
	r := gin.New()
	r.Use(LoggerMiddleware(zap.NewNop())) // Using No-op logger for this test

	// Define a test endpoint
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Send the request
	r.ServeHTTP(rr, req)

	// Assert the status code and body response
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test", rr.Body.String())
}

// TestLoggerMiddleware_LogRequest checks the log output of the LoggerMiddleware
func TestLoggerMiddleware_LogRequest(t *testing.T) {
	// Set up a buffer to capture the log output
	var buf bytes.Buffer

	// Set up the logger to write to the buffer
	loggerConfig := zap.NewProductionEncoderConfig()
	loggerConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(loggerConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(&buf), zapcore.InfoLevel)
	logger := zap.New(core)

	// Create a Gin router and add LoggerMiddleware using the custom logger
	r := gin.New()
	r.Use(LoggerMiddleware(logger)) // Pass the custom logger to the middleware

	// Define a test endpoint
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "test")
	})

	// Create a new HTTP request
	req, _ := http.NewRequest("GET", "/test", nil)

	// Create a recorder to capture the response
	rr := httptest.NewRecorder()

	// Send the request
	r.ServeHTTP(rr, req)

	// Check if the status code and response body are correct
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test", rr.Body.String())

	// Assert the log contains the expected information
	logOutput := buf.String()

	// Assert that the log contains the method and URI
	assert.Contains(t, logOutput, "\"method\":\"GET\"")
	assert.Contains(t, logOutput, "\"uri\":\"/test\"")

	// Check if the status code and content length are present in the log
	assert.Contains(t, logOutput, "\"status_code\":200")

}
