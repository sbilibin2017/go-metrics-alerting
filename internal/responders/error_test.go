package responders

import (
	"errors"
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// performRequest executes an HTTP request and returns the response recorder
func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestErrorResponder_Respond(t *testing.T) {
	// Set Gin mode to TestMode
	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Test 1: Simulate APIError response
	r.GET("/test", func(c *gin.Context) {
		responder := &ErrorResponder{C: c}
		apiErr := &types.APIErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Bad Request",
		}
		responder.Respond(apiErr)
	})

	t.Run("Should return correct APIError message", func(t *testing.T) {
		// Send a request to the route
		w := performRequest(r, "GET", "/test")

		// Verify the HTTP status code
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Verify the response body
		assert.Equal(t, "Bad Request", w.Body.String())
	})

	// Test 2: Simulate an unknown error (using the default Internal Server Error)
	r.GET("/test-default-error", func(c *gin.Context) {
		responder := &ErrorResponder{C: c}
		var _ error = errors.New("some unknown error") // Simulate an unknown error
		responder.Respond(nil)                         // Pass nil to trigger the internal server error
	})

	t.Run("Should return Internal Server Error for unknown error", func(t *testing.T) {
		// Send a request to the route
		w := performRequest(r, "GET", "/test-default-error")

		// Verify the HTTP status code
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verify the response body
		assert.Equal(t, "Internal Server Error", w.Body.String())
	})
}
