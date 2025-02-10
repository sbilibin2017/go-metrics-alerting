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

// Мокаем интерфейс UpdateValueService
type MockUpdateValueService struct {
	mock.Mock
}

func (m *MockUpdateValueService) UpdateMetricValue(ctx context.Context, req *types.UpdateMetricValueRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// Мокаем интерфейсы валидаторов
type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(value string) error {
	args := m.Called(value)
	return args.Error(0)
}

func setupMocks() (*MockUpdateValueService, *MockValidator, *MockValidator, *MockValidator, *MockValidator, *MockValidator) {
	// Создаем моки для сервисов и валидаторов
	mockService := new(MockUpdateValueService)
	mockMetricTypeValidator := new(MockValidator)
	mockMetricNameValidator := new(MockValidator)
	mockMetricValueValidator := new(MockValidator)
	mockGaugeValueValidator := new(MockValidator)
	mockCounterValueValidator := new(MockValidator)

	return mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator
}

func TestValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим тестируемые моки
	mockMetricTypeValidator.On("Validate", "gauge").Return(nil)
	mockMetricNameValidator.On("Validate", "metric1").Return(nil)
	mockMetricValueValidator.On("Validate", "10.5").Return(nil)
	mockGaugeValueValidator.On("Validate", "10.5").Return(nil)
	mockService.On("UpdateMetricValue", mock.Anything, mock.Anything).Return(nil)

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/metric1/10.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и ответ
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Metric updated", w.Body.String())
}

func TestInvalidMetricType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим тестируемые моки
	mockMetricTypeValidator.On("Validate", "invalid_type").Return(errors.New("Invalid metric type"))

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	req, _ := http.NewRequest(http.MethodPost, "/update/invalid_type/metric1/10.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и ответ
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid metric type", w.Body.String())
}

func TestInvalidMetricName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим тестируемые моки
	mockMetricTypeValidator.On("Validate", "gauge").Return(nil)
	mockMetricNameValidator.On("Validate", "invalid_metric").Return(errors.New("Invalid metric name"))

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/invalid_metric/10.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и ответ
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Invalid metric name", w.Body.String())
}

func TestInvalidMetricValue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим тестируемые моки
	mockMetricTypeValidator.On("Validate", "gauge").Return(nil)
	mockMetricNameValidator.On("Validate", "metric1").Return(nil)
	mockMetricValueValidator.On("Validate", "invalid_value").Return(errors.New("Invalid metric value"))

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/metric1/invalid_value", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и ответ
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid metric value", w.Body.String())
}

// Тест на ошибку сервиса
func TestServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Настроим моки
	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим, что будет происходить при вызове метода UpdateMetricValue
	mockMetricTypeValidator.On("Validate", "gauge").Return(nil)
	mockMetricNameValidator.On("Validate", "metric1").Return(nil)
	mockMetricValueValidator.On("Validate", "10.5").Return(nil)
	mockGaugeValueValidator.On("Validate", "10.5").Return(nil)
	mockService.On("UpdateMetricValue", mock.Anything, mock.Anything).Return(&apierror.APIError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
	})

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	// Выполнение запроса
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/metric1/10.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}

func TestInternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Настроим моки
	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим, что будет происходить при вызове метода UpdateMetricValue
	mockMetricTypeValidator.On("Validate", "gauge").Return(nil)
	mockMetricNameValidator.On("Validate", "metric1").Return(nil)
	mockMetricValueValidator.On("Validate", "10.5").Return(nil)
	mockGaugeValueValidator.On("Validate", "10.5").Return(nil)

	// Настроим, что сервис вызывает обычную ошибку (не APIError)
	mockService.On("UpdateMetricValue", mock.Anything, mock.Anything).Return(errors.New("Unexpected error"))

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	// Выполнение запроса
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/metric1/10.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}

func TestGaugeValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Настроим моки
	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим, что будет происходить при вызове метода UpdateMetricValue
	mockMetricTypeValidator.On("Validate", "gauge").Return(nil)
	mockMetricNameValidator.On("Validate", "metric1").Return(nil)
	mockMetricValueValidator.On("Validate", "10.5").Return(nil)

	// Настроим, что валидация для Gauge возвращает ошибку
	mockGaugeValueValidator.On("Validate", "10.5").Return(errors.New("Invalid Gauge value"))

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	// Выполнение запроса
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/metric1/10.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid Gauge value", w.Body.String())
}

func TestCounterValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Настроим моки
	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим, что будет происходить при вызове метода UpdateMetricValue
	mockMetricTypeValidator.On("Validate", "counter").Return(nil)
	mockMetricNameValidator.On("Validate", "metric1").Return(nil)
	mockMetricValueValidator.On("Validate", "10.5").Return(nil)

	// Настроим, что валидация для Counter возвращает ошибку
	mockCounterValueValidator.On("Validate", "10.5").Return(errors.New("Invalid Counter value"))

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	// Выполнение запроса
	req, _ := http.NewRequest(http.MethodPost, "/update/counter/metric1/10.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Invalid Counter value", w.Body.String())
}

func TestNoRouteHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Настройка моки для всех валидаторов (можно оставить как заглушки)
	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	// Выполнение запроса на несуществующий маршрут
	req, _ := http.NewRequest(http.MethodGet, "/nonexistent-route", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка, что сервер вернул код 404 и сообщение "Route not found"
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Route not found", w.Body.String())
}

func TestUnknownMetricType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Настроим моки
	mockService, mockMetricTypeValidator, mockMetricNameValidator, mockMetricValueValidator, mockGaugeValueValidator, mockCounterValueValidator := setupMocks()

	// Настроим, что будет происходить при вызове метода Validate
	mockMetricTypeValidator.On("Validate", "unknown_type").Return(errors.New("Unsupported metric type"))

	// Настройка роутинга
	router := gin.Default()
	RegisterUpdateMetricValueHandler(router, mockService,
		mockMetricTypeValidator,
		mockMetricNameValidator,
		mockMetricValueValidator,
		mockGaugeValueValidator,
		mockCounterValueValidator,
	)

	// Выполнение запроса с неизвестным типом метрики
	req, _ := http.NewRequest(http.MethodPost, "/update/unknown_type/metric1/10.5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверка ответа
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Unsupported metric type", w.Body.String()) // Сообщение об ошибке
}
