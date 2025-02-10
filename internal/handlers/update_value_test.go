package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
)

// Mock сервис

type MockUpdateValueService struct {
	mock.Mock
}

func (m *MockUpdateValueService) UpdateMetricValue(ctx context.Context, req *types.UpdateMetricValueRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// Mock валидаторов

type MockValidator struct {
	mock.Mock
}

func (m *MockValidator) Validate(value string) error {
	args := m.Called(value)
	return args.Error(0)
}

func setupHandler() (*UpdateValueHandler, *gin.Engine, *MockUpdateValueService, *MockValidator, *MockValidator, *MockValidator, *MockValidator, *MockValidator) {
	gin.SetMode(gin.TestMode)

	service := new(MockUpdateValueService)
	typeValidator := new(MockValidator)
	nameValidator := new(MockValidator)
	valueValidator := new(MockValidator)
	gaugeValidator := new(MockValidator)
	counterValidator := new(MockValidator)

	h := &UpdateValueHandler{
		service:               service,
		metricTypeValidator:   typeValidator,
		metricNameValidator:   nameValidator,
		metricValueValidator:  valueValidator,
		gaugeValueValidator:   gaugeValidator,
		counterValueValidator: counterValidator,
	}

	r := gin.Default()
	h.RegisterRoutes(r)

	return h, r, service, typeValidator, nameValidator, valueValidator, gaugeValidator, counterValidator
}

func TestSuccessfulMetricUpdate(t *testing.T) {
	_, r, service, typeValidator, nameValidator, valueValidator, gaugeValidator, _ := setupHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/cpu/42", nil)

	typeValidator.On("Validate", "gauge").Return(nil)
	nameValidator.On("Validate", "cpu").Return(nil)
	valueValidator.On("Validate", "42").Return(nil)
	gaugeValidator.On("Validate", "42").Return(nil)
	service.On("UpdateMetricValue", mock.Anything, mock.Anything).Return(nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Metric updated", w.Body.String())
}

func TestInvalidMetricType(t *testing.T) {
	_, r, _, typeValidator, _, _, _, _ := setupHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/update/invalid/cpu/42", nil)

	typeValidator.On("Validate", "invalid").Return(errors.New("invalid type"))

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid type", w.Body.String())
}

func TestServiceError(t *testing.T) {
	_, r, service, typeValidator, nameValidator, valueValidator, gaugeValidator, _ := setupHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/cpu/42", nil)

	typeValidator.On("Validate", "gauge").Return(nil)
	nameValidator.On("Validate", "cpu").Return(nil)
	valueValidator.On("Validate", "42").Return(nil)
	gaugeValidator.On("Validate", "42").Return(nil)
	service.On("UpdateMetricValue", mock.Anything, mock.Anything).Return(&apierror.APIError{Code: http.StatusInternalServerError, Message: "error updating metric"})

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "error updating metric", w.Body.String())
}

func TestMetricNameValidationError(t *testing.T) {
	// Настройка теста
	_, r, _, typeValidator, nameValidator, valueValidator, gaugeValidator, _ := setupHandler()

	// Мокаем ошибку для имени метрики
	nameValidator.On("Validate", "cpu").Return(errors.New("invalid metric name"))

	// Мокаем успешную валидацию для других валидаторов
	typeValidator.On("Validate", "gauge").Return(nil)
	valueValidator.On("Validate", "42").Return(nil)
	gaugeValidator.On("Validate", "42").Return(nil)

	// Выполняем HTTP запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/cpu/42", nil)

	// Запускаем запрос через маршруты
	r.ServeHTTP(w, req)

	// Проверяем, что мы получаем код ошибки 404 и правильное сообщение
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "invalid metric name", w.Body.String())
}

func TestMetricValueValidationError(t *testing.T) {
	// Настройка теста
	_, r, _, typeValidator, nameValidator, valueValidator, gaugeValidator, _ := setupHandler()

	// Мокаем успешную валидацию для типа метрики и имени метрики
	typeValidator.On("Validate", "gauge").Return(nil)
	nameValidator.On("Validate", "cpu").Return(nil)

	// Мокаем ошибку валидации для значения метрики
	valueValidator.On("Validate", "invalid_value").Return(errors.New("invalid metric value"))
	gaugeValidator.On("Validate", "invalid_value").Return(nil)

	// Выполняем HTTP запрос
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/cpu/invalid_value", nil)

	// Запускаем запрос через маршруты
	r.ServeHTTP(w, req)

	// Проверяем, что мы получаем код ошибки 400 и правильное сообщение
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid metric value", w.Body.String())
}

func TestGaugeValueValidationError(t *testing.T) {
	_, r, _, typeValidator, nameValidator, valueValidator, gaugeValidator, _ := setupHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/cpu/invalid_value", nil)

	// Мокаем ошибку валидации для типа метрики "gauge"
	typeValidator.On("Validate", "gauge").Return(nil)
	nameValidator.On("Validate", "cpu").Return(nil)
	valueValidator.On("Validate", "invalid_value").Return(nil)
	gaugeValidator.On("Validate", "invalid_value").Return(errors.New("invalid gauge value"))

	r.ServeHTTP(w, req)

	// Проверка на ошибку с кодом 400 и сообщением ошибки
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid gauge value", w.Body.String())
}

func TestCounterValueValidationError(t *testing.T) {
	_, r, _, typeValidator, nameValidator, valueValidator, _, counterValidator := setupHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/update/counter/cpu/42", nil)

	// Мокаем ошибку валидации для типа метрики "counter"
	typeValidator.On("Validate", "counter").Return(nil)
	nameValidator.On("Validate", "cpu").Return(nil)
	valueValidator.On("Validate", "42").Return(nil)
	counterValidator.On("Validate", "42").Return(errors.New("invalid counter value"))

	r.ServeHTTP(w, req)

	// Проверка на ошибку с кодом 400 и сообщением ошибки
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "invalid counter value", w.Body.String())
}

func TestNoRouteHandler(t *testing.T) {
	_, r, _, _, _, _, _, _ := setupHandler()

	// Создаём запрос на несуществующий маршрут
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/nonexistent/route", nil)

	// Запускаем запрос через маршруты
	r.ServeHTTP(w, req)

	// Проверяем, что мы получаем код ошибки 404 и правильное сообщение
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Route not found", w.Body.String())
}

func TestServiceInternalError(t *testing.T) {
	_, r, service, typeValidator, nameValidator, valueValidator, gaugeValidator, _ := setupHandler()

	// Создаём запрос для обновления метрики
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/update/gauge/cpu/42", nil)

	// Настроим моки для валидаторов
	typeValidator.On("Validate", "gauge").Return(nil)
	nameValidator.On("Validate", "cpu").Return(nil)
	valueValidator.On("Validate", "42").Return(nil)
	gaugeValidator.On("Validate", "42").Return(nil)

	// Настроим мок сервиса, чтобы он возвращал ошибку
	service.On("UpdateMetricValue", mock.Anything, mock.Anything).Return(errors.New("service error"))

	// Запускаем запрос
	r.ServeHTTP(w, req)

	// Проверяем, что получен статус 500 и правильное сообщение
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}
