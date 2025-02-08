package handlers

import (
	"errors"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для сервиса GetValueService
type MockGetValueService struct {
	mock.Mock
}

func (m *MockGetValueService) GetMetricValue(req *GetMetricValueRequest) (string, error) {
	args := m.Called(req)
	return args.String(0), args.Error(1)
}

func TestRegisterGetMetricValueHandler_ValidationError_MissingMetricType(t *testing.T) {
	mockService := new(MockGetValueService)
	router := gin.Default()
	RegisterGetMetricValueHandler(router, mockService)

	req, err := http.NewRequest("GET", "/value//cpu", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Metric type is required", w.Body.String())
}

func TestRegisterGetMetricValueHandler_ValidationError_MissingMetricName(t *testing.T) {
	mockService := new(MockGetValueService)
	router := gin.Default()
	RegisterGetMetricValueHandler(router, mockService)

	req, err := http.NewRequest("GET", "/value/gauge", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Metric name is required", w.Body.String())
}

func TestRegisterGetMetricValueHandler_Success(t *testing.T) {
	mockService := new(MockGetValueService)
	router := gin.Default()
	RegisterGetMetricValueHandler(router, mockService)

	mockService.On("GetMetricValue", &GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("100", nil)

	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "100", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Header().Get("Date"))
	assert.Equal(t, "3", w.Header().Get("Content-Length"))
}

func TestRegisterGetMetricValueHandler_ServiceError_NotFound(t *testing.T) {
	mockService := new(MockGetValueService)
	router := gin.Default()
	RegisterGetMetricValueHandler(router, mockService)

	mockService.On("GetMetricValue", &GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("", &apierror.APIError{
		Code:    http.StatusNotFound,
		Message: "Metric not found",
	})

	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Metric not found", w.Body.String())
}

func TestRegisterGetMetricValueHandler_ServiceError_Internal(t *testing.T) {
	mockService := new(MockGetValueService)
	router := gin.Default()
	RegisterGetMetricValueHandler(router, mockService)

	mockService.On("GetMetricValue", &GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("", errors.New("internal error"))

	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}
