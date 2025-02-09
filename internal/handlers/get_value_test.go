package handlers

import (
	"context"
	"errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок-сервис для GetValueService
type MockGetValueService struct {
	mock.Mock
}

func (m *MockGetValueService) GetMetricValue(ctx context.Context, req *types.GetMetricValueRequest) (string, error) {
	args := m.Called(ctx, req)
	return args.String(0), args.Error(1)
}

func setupTestRouter(svc GetValueService) *gin.Engine {
	r := gin.Default()
	RegisterGetMetricValueHandler(r, svc)
	return r
}

func TestGetMetricValueHandler_Success(t *testing.T) {
	// Arrange
	mockService := new(MockGetValueService)
	router := setupTestRouter(mockService)

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("100", nil)

	// Создаем запрос
	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "100", w.Body.String())
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Header().Get("Date"))

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}

func TestGetMetricValueHandler_NotFound(t *testing.T) {
	// Arrange
	mockService := new(MockGetValueService)
	router := setupTestRouter(mockService)

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("", &apierror.APIError{
		Code:    http.StatusNotFound,
		Message: "Metric not found",
	})

	// Создаем запрос
	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Metric not found", w.Body.String())

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}

func TestGetMetricValueHandler_InternalServerError(t *testing.T) {
	// Arrange
	mockService := new(MockGetValueService)
	router := setupTestRouter(mockService)

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("", errors.New("internal error"))

	// Создаем запрос
	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}

func TestGetMetricValueByTypeHandler_Success(t *testing.T) {
	// Arrange
	mockService := new(MockGetValueService)
	router := setupTestRouter(mockService)

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: "counter",
	}).Return("200", nil)

	// Создаем запрос
	req, err := http.NewRequest("GET", "/value/counter", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "200", w.Body.String())

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}

func TestGetMetricValueByTypeHandler_InternalServerError(t *testing.T) {
	// Arrange
	mockService := new(MockGetValueService)
	router := setupTestRouter(mockService)

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: "error",
	}).Return("", errors.New("internal error"))

	// Создаем запрос
	req, err := http.NewRequest("GET", "/value/error", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}

func TestGetMetricValueByTypeHandler_ServiceError(t *testing.T) {
	// Arrange
	mockService := new(MockGetValueService)
	router := gin.Default()
	router.GET("/value/:type", func(c *gin.Context) {
		getMetricValueByTypeHandler(mockService, c)
	})

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: "gauge",
	}).Return("", errors.New("internal error"))

	// Создаем запрос
	req, err := http.NewRequest("GET", "/value/gauge", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}

func TestGetMetricValueByTypeHandler_APINotFoundError(t *testing.T) {
	// Arrange
	mockService := new(MockGetValueService)
	router := gin.Default()
	router.GET("/value/:type", func(c *gin.Context) {
		getMetricValueByTypeHandler(mockService, c)
	})

	// Устанавливаем ожидание на мок-сервис
	mockService.On("GetMetricValue", mock.Anything, &types.GetMetricValueRequest{
		Type: "gauge",
	}).Return("", &apierror.APIError{
		Code:    http.StatusNotFound,
		Message: "Metric type not found",
	})

	// Создаем запрос
	req, err := http.NewRequest("GET", "/value/gauge", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем результат
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Metric type not found", w.Body.String())

	// Проверка того, что метод был вызван
	mockService.AssertExpectations(t)
}
