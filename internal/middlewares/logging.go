package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Global logger instance
var logger *zap.Logger

// Initialize zap logger in the init function
func init() {
	// Create a production-level zap logger
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic("failed to initialize zap logger")
	}

	// Defer flushing logs on program termination
	defer logger.Sync()
}

// LoggingMiddleware logs incoming HTTP requests and their corresponding responses
func LoggingMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Create a custom response writer to capture status and size
			responseRecorder := &responseWriter{ResponseWriter: w}

			// Capture start time for request duration
			start := time.Now()

			// Log the incoming request using zap logger
			logger.Info("Request received",
				zap.String("method", r.Method),
				zap.String("uri", r.RequestURI),
				zap.String("host", r.Host),
			)

			// Serve the HTTP request
			next.ServeHTTP(responseRecorder, r)

			// Calculate request duration
			duration := time.Since(start)

			// Log the response details using zap logger
			logger.Info("Response sent",
				zap.Int("status", responseRecorder.StatusCode()),
				zap.Int("size", responseRecorder.Size()),
				zap.Duration("duration", duration),
			)
		})
	}
}

// Custom response writer to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(p []byte) (n int, err error) {
	n, err = rw.ResponseWriter.Write(p)
	rw.size += n
	return n, err
}

func (rw *responseWriter) StatusCode() int {
	return rw.statusCode
}

func (rw *responseWriter) Size() int {
	return rw.size
}
