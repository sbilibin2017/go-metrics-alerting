package handlers

import (
	"context"
	"fmt"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type mockMetricService struct {
	UpdateErr error
	GetResp   string
	GetErr    error
	ListResp  []*types.MetricResponse
}

func (m *mockMetricService) UpdateMetric(ctx context.Context, req *types.UpdateMetricValueRequest) error {
	return m.UpdateErr
}

func (m *mockMetricService) GetMetric(ctx context.Context, req *types.GetMetricValueRequest) (string, error) {
	return m.GetResp, m.GetErr
}

func (m *mockMetricService) ListMetrics(ctx context.Context) []*types.MetricResponse {
	return m.ListResp
}

func TestUpdateMetricHandler(t *testing.T) {
	service := &mockMetricService{}
	r := gin.Default()
	r.POST("/update/:type/:name/:value", UpdateMetricHandler(service))

	req, _ := http.NewRequest("POST", "/update/gauge/cpu/0.5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Metric updated", w.Body.String())
}

func TestGetMetricHandler(t *testing.T) {
	service := &mockMetricService{GetResp: "42"}
	r := gin.Default()
	r.GET("/value/:type/:name", GetMetricHandler(service))

	req, _ := http.NewRequest("GET", "/value/:type/:name", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "42", w.Body.String())
}

func TestListMetricsHandler_WithMetrics(t *testing.T) {
	// Test with a non-empty list of metrics
	service := &mockMetricService{ListResp: []*types.MetricResponse{
		{Name: "cpu", Value: "0.5"},
	}}
	r := gin.Default()
	r.GET("/", ListMetricsHandler(service))

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<li>cpu: 0.5</li>")
}

func TestListMetricsHandler_EmptyList(t *testing.T) {
	// Test with an empty list of metrics
	service := &mockMetricService{ListResp: []*types.MetricResponse{}}
	r := gin.Default()
	r.GET("/", ListMetricsHandler(service))

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<li>No metrics available</li>")
}

// Test handleError for APIError (custom error type)
func TestHandleError_APIError(t *testing.T) {
	r := gin.Default()

	// Create an instance of the APIError with a custom message and code
	apiErr := &apierror.APIError{
		Code:    http.StatusBadRequest,
		Message: "Custom API error occurred",
	}

	// Create a new Gin context (we'll mock this)
	r.GET("/error", func(c *gin.Context) {
		handleError(c, apiErr) // Simulating error handler
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the correct status code and error message are returned
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Custom API error occurred", w.Body.String())
}

// Test handleError for generic error (non-APIError)
func TestHandleError_GenericError(t *testing.T) {
	r := gin.Default()

	// Create a generic error
	genericErr := fmt.Errorf("a generic error occurred")

	// Create a new Gin context (we'll mock this)
	r.GET("/error", func(c *gin.Context) {
		handleError(c, genericErr) // Simulating error handler
	})

	req, _ := http.NewRequest("GET", "/error", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the correct status code and error message are returned
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}

func TestUpdateMetricHandler_Success(t *testing.T) {
	service := &mockMetricService{}
	r := gin.Default()
	r.POST("/update/:type/:name/:value", UpdateMetricHandler(service))

	req, _ := http.NewRequest("POST", "/update/gauge/cpu/0.5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Metric updated", w.Body.String())
}

func TestUpdateMetricHandler_APIError(t *testing.T) {
	// Simulate an APIError from the service
	service := &mockMetricService{
		UpdateErr: &apierror.APIError{
			Code:    http.StatusBadRequest,
			Message: "Invalid metric value",
		},
	}

	r := gin.Default()
	r.POST("/update/:type/:name/:value", UpdateMetricHandler(service))

	req, _ := http.NewRequest("POST", "/update/gauge/cpu/invalid-value", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the response contains the APIError's code and message
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid metric value", w.Body.String())
}

func TestUpdateMetricHandler_GenericError(t *testing.T) {
	// Simulate a generic error from the service
	service := &mockMetricService{
		UpdateErr: fmt.Errorf("some internal error occurred"),
	}

	r := gin.Default()
	r.POST("/update/:type/:name/:value", UpdateMetricHandler(service))

	req, _ := http.NewRequest("POST", "/update/gauge/cpu/0.5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that a generic internal server error is returned
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}

func TestGetMetricHandler_Success(t *testing.T) {
	// Simulate the GetMetric response
	service := &mockMetricService{GetResp: "42"}
	r := gin.Default()
	r.GET("/value/:type/:name", GetMetricHandler(service))

	req, _ := http.NewRequest("GET", "/value/gauge/cpu", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the response contains the metric value
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "42", w.Body.String())
}

func TestGetMetricHandler_APIError(t *testing.T) {
	// Simulate an APIError from the service
	service := &mockMetricService{
		GetErr: &apierror.APIError{
			Code:    http.StatusNotFound,
			Message: "Metric not found",
		},
	}

	r := gin.Default()
	r.GET("/value/:type/:name", GetMetricHandler(service))

	req, _ := http.NewRequest("GET", "/value/gauge/cpu", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that the response contains the APIError's code and message
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Metric not found", w.Body.String())
}

func TestGetMetricHandler_GenericError(t *testing.T) {
	// Simulate a generic error from the service
	service := &mockMetricService{
		GetErr: fmt.Errorf("some internal error occurred"),
	}

	r := gin.Default()
	r.GET("/value/:type/:name", GetMetricHandler(service))

	req, _ := http.NewRequest("GET", "/value/gauge/cpu", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Check that a generic internal server error is returned
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}
