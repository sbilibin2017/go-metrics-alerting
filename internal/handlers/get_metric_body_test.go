package handlers

import (
	"bytes"
	"go-metrics-alerting/internal/domain"

	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the service to be used in the handler test
type MockGetMetricBodyService struct {
	mock.Mock
}

func (m *MockGetMetricBodyService) GetMetric(id string, mtype domain.MetricType) (*domain.Metrics, error) {
	args := m.Called(id, mtype)
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Metrics), args.Error(1)
}

func TestGetMetricBodyHandler_Success(t *testing.T) {
	mockService := new(MockGetMetricBodyService)

	// Prepare a valid metric response
	metric := &domain.Metrics{
		ID:    "metric1",
		MType: domain.Counter,
		Delta: new(int64),
		Value: new(float64),
	}
	*metric.Delta = 10
	*metric.Value = 15.5

	// Mock the service method
	mockService.On("GetMetric", "metric1", domain.Counter).Return(metric, nil)

	// Prepare the valid request body
	reqBody := types.GetMetricRequest{
		ID:    "metric1",
		MType: "counter",
	}
	data, _ := easyjson.Marshal(reqBody)

	r := httptest.NewRequest(http.MethodPost, "/metrics", bytes.NewReader(data))
	w := httptest.NewRecorder()

	// Create and serve the handler
	handler := GetMetricBodyHandler(mockService)
	handler.ServeHTTP(w, r)

	// Check the response
	assert.Equal(t, http.StatusOK, w.Code)

	// Expected response
	expectedResponse := types.GetMetricBodyResponse{
		GetMetricRequest: types.GetMetricRequest{
			ID:    "metric1",
			MType: "counter",
		},
		Delta: metric.Delta,
		Value: metric.Value,
	}

	// Unmarshal the response body to check the content
	var resp types.GetMetricBodyResponse
	err := easyjson.Unmarshal(w.Body.Bytes(), &resp)
	assert.Nil(t, err)
	assert.Equal(t, expectedResponse, resp)

	mockService.AssertExpectations(t)
}

func TestGetMetricBodyHandler_ValidationError(t *testing.T) {
	// Create a request with an empty ID
	reqBody := types.GetMetricRequest{
		ID:    "",
		MType: "counter",
	}
	data, _ := easyjson.Marshal(reqBody)

	r := httptest.NewRequest(http.MethodPost, "/metrics", bytes.NewReader(data))
	w := httptest.NewRecorder()

	mockService := new(MockGetMetricBodyService)
	handler := GetMetricBodyHandler(mockService)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetMetricBodyHandler_ServiceError(t *testing.T) {
	// Prepare a valid request body
	reqBody := types.GetMetricRequest{
		ID:    "metric1",
		MType: "counter",
	}
	data, _ := easyjson.Marshal(reqBody)

	r := httptest.NewRequest(http.MethodPost, "/metrics", bytes.NewReader(data))
	w := httptest.NewRecorder()

	mockService := new(MockGetMetricBodyService)

	// Mock a service error
	mockService.On("GetMetric", "metric1", domain.Counter).Return(nil, assert.AnError)

	handler := GetMetricBodyHandler(mockService)
	handler.ServeHTTP(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetMetricBodyHandler_InvalidJSONBody(t *testing.T) {
	// Create an invalid JSON body (e.g., missing closing brace or malformed structure)
	invalidJSON := []byte(`{"id": "metric1", "mtype": "counter"`) // Missing closing brace

	// Create a request with the invalid JSON body
	r := httptest.NewRequest(http.MethodPost, "/metrics", bytes.NewReader(invalidJSON))
	w := httptest.NewRecorder()

	// Create the mock service
	mockService := new(MockGetMetricBodyService)

	// Create the handler
	handler := GetMetricBodyHandler(mockService)

	// Call the handler
	handler.ServeHTTP(w, r)

	// Check the response status and body
	assert.Equal(t, http.StatusBadRequest, w.Code)

}
