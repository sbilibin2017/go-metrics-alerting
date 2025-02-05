package handlers

import (
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUpdateValueService struct {
	mock.Mock
}

func (m *MockUpdateValueService) UpdateMetricValue(req *UpdateMetricValueRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func TestRegisterUpdateValueHandler_Success(t *testing.T) {
	mockService := new(MockUpdateValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterUpdateValueHandler(router, mockService)

	// Мокируем успешное обновление метрики
	mockService.On("UpdateMetricValue", &UpdateMetricValueRequest{
		Type:  "gauge",
		Name:  "cpu",
		Value: "99",
	}).Return(nil)

	// Отправляем запрос
	req, err := http.NewRequest("POST", "/update/gauge/cpu/99", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и тело ответа
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Metric updated", w.Body.String())

	mockService.AssertExpectations(t)
}

func TestRegisterUpdateValueHandler_ValidationError_EmptyMetricValue(t *testing.T) {
	mockService := new(MockUpdateValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterUpdateValueHandler(router, mockService)

	// Отправляем запрос с пустым значением метрики
	req, err := http.NewRequest("POST", "/update/gauge/cpu/", nil) // Пустое значение метрики
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и тело ответа (ошибка валидации значения)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "404 page not found", w.Body.String())

	mockService.AssertExpectations(t)
}

func TestRegisterUpdateValueHandler_ValidationError_MetricType(t *testing.T) {
	mockService := new(MockUpdateValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterUpdateValueHandler(router, mockService)

	// Отправляем запрос с пустым типом метрики (обратите внимание на два слэша)
	req, err := http.NewRequest("POST", "/update//cpu/99", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и тело ответа (ошибка валидации типа)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Metric type is required", w.Body.String())
}

func TestRegisterUpdateValueHandler_ValidationError_MetricName(t *testing.T) {
	mockService := new(MockUpdateValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterUpdateValueHandler(router, mockService)

	// Отправляем запрос с пустым именем метрики
	req, err := http.NewRequest("POST", "/update/gauge//99", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и тело ответа (ошибка валидации имени)
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Metric name is required", w.Body.String())

	mockService.AssertExpectations(t)
}

func TestRegisterUpdateValueHandler_InternalServerError(t *testing.T) {
	mockService := new(MockUpdateValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterUpdateValueHandler(router, mockService)

	// Мокируем ошибку при обновлении метрики
	mockService.On("UpdateMetricValue", &UpdateMetricValueRequest{
		Type:  "gauge",
		Name:  "cpu",
		Value: "99",
	}).Return(&apierror.APIError{
		Code:    http.StatusInternalServerError,
		Message: "Internal error",
	})

	// Отправляем запрос
	req, err := http.NewRequest("POST", "/update/gauge/cpu/99", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем статус и тело ответа (внутренняя ошибка)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())

	mockService.AssertExpectations(t)
}
