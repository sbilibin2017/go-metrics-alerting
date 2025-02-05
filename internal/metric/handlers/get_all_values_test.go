package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Mock-сервис, который возвращает метрики
type mockGetAllValuesService struct {
	metrics []*MetricResponse
}

func (m *mockGetAllValuesService) GetAllMetricValues() []*MetricResponse {
	return m.metrics
}

func TestRegisterGetAllMetricValuesHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &mockGetAllValuesService{
		metrics: []*MetricResponse{
			{Type: "gauge", Name: "cpu", Value: "95.5"},
			{Type: "counter", Name: "hits", Value: "200"},
		},
	}

	RegisterGetAllMetricValuesHandler(router, mockService)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<h1>Metrics List</h1>")
	assert.Contains(t, w.Body.String(), "cpu: 95.5")
	assert.Contains(t, w.Body.String(), "hits: 200")
}

func TestRegisterGetAllMetricValuesHandler_EmptyList(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := &mockGetAllValuesService{
		metrics: []*MetricResponse{},
	}

	RegisterGetAllMetricValuesHandler(router, mockService)

	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<h1>Metrics List</h1>")
	assert.NotContains(t, w.Body.String(), "<li>")
}
