package handlers_test

import (
	"bytes"
	"encoding/json"
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/handlers"
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock service for UpdateMetricsService
type MockUpdateMetricsService struct {
	mock.Mock
}

func (m *MockUpdateMetricsService) UpdateMetricValue(metric *domain.Metric) (*domain.Metric, error) {
	args := m.Called(metric)
	return args.Get(0).(*domain.Metric), args.Error(1)
}

func TestUpdateMetricsBodyHandler(t *testing.T) {
	// Setup mock service
	mockService := new(MockUpdateMetricsService)
	mockService.On("UpdateMetricValue", mock.AnythingOfType("*domain.Metric")).Return(&domain.Metric{
		ID:    "metric1",
		MType: domain.Counter,
		Value: "123.45",
	}, nil)

	// Initialize Gin router and register route
	r := gin.Default()
	r.POST("/update/", handlers.UpdateMetricsBodyHandler(mockService))

	// Create a test request with JSON body
	requestBody := types.UpdateMetricBodyRequest{
		ID:    "metric1",
		MType: "counter",
		// Only using Value, not Delta
		Value: float64Ptr(123.45),
	}
	jsonData, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	recorder := httptest.NewRecorder()
	r.ServeHTTP(recorder, req)

	// Assertions
	assert.Equal(t, http.StatusOK, recorder.Code)
	var response map[string]interface{}
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "metric1", response["id"])
	assert.Equal(t, "Counter", response["mtype"])
	assert.Equal(t, 123.45, response["value"])

	// Verify the mock service was called
	mockService.AssertExpectations(t)
}

func float64Ptr(v float64) *float64 {
	return &v
}
