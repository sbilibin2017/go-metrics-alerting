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

func TestRegisterGetMetricValueHandler_ValidationError_EmptyMetricName(t *testing.T) {
	// Мокируем сервис
	mockService := new(MockGetValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterGetMetricValueHandler(router, mockService)

	// Отправляем запрос с пустым именем метрики
	req, err := http.NewRequest("GET", "/value/gauge/", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа и тело
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "404 page not found", w.Body.String())
}

func TestRegisterGetMetricValueHandler_ValidationError_MetricType(t *testing.T) {
	// Мокируем сервис
	mockService := new(MockGetValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterGetMetricValueHandler(router, mockService)

	// Отправляем запрос с пустым типом метрики
	req, err := http.NewRequest("GET", "/value//cpu", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа и тело
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "Metric type is required", w.Body.String())
}

func TestRegisterGetMetricValueHandler_ValidationError_MetricName(t *testing.T) {
	// Создаем тестовый запрос с отсутствующим параметром 'name'
	router := gin.Default()
	handler := &MockGetValueService{}
	RegisterGetMetricValueHandler(router, handler)

	req, _ := http.NewRequest("GET", "/value/gauge/", nil)
	w := httptest.NewRecorder()

	// Выполняем запрос
	router.ServeHTTP(w, req)

	// Проверяем, что возвращается правильный код ошибки и сообщение
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "404 page not found", w.Body.String())
}

func TestRegisterGetMetricValueHandler_Success(t *testing.T) {
	// Мокируем сервис
	mockService := new(MockGetValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterGetMetricValueHandler(router, mockService)

	// Ожидаем, что сервис вернет значение "100"
	mockService.On("GetMetricValue", &GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("100", nil)

	// Отправляем корректный запрос
	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа и тело
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "100", w.Body.String())

	// Проверяем заголовки
	assert.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
	assert.NotEmpty(t, w.Header().Get("Date"))
	assert.Equal(t, "3", w.Header().Get("Content-Length"))
}

func TestRegisterGetMetricValueHandler_ServiceError(t *testing.T) {
	// Мокируем сервис
	mockService := new(MockGetValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterGetMetricValueHandler(router, mockService)

	// Ожидаем, что сервис вернет ошибку с кодом 404
	mockService.On("GetMetricValue", &GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("", &apierror.APIError{
		Code:    http.StatusNotFound,
		Message: "Metric not found",
	})

	// Отправляем запрос с меткой, которой нет в хранилище
	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа и тело
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Metric not found", w.Body.String())
}

func TestRegisterGetMetricValueHandler_InternalServerError(t *testing.T) {
	// Мокируем сервис
	mockService := new(MockGetValueService)
	router := gin.Default()

	// Регистрация обработчика
	RegisterGetMetricValueHandler(router, mockService)

	// Ожидаем, что сервис вернет ошибку
	mockService.On("GetMetricValue", &GetMetricValueRequest{
		Type: "gauge",
		Name: "cpu",
	}).Return("", errors.New("internal error"))

	// Отправляем запрос с ошибкой
	req, err := http.NewRequest("GET", "/value/gauge/cpu", nil)
	assert.NoError(t, err)

	// Запускаем тестирование
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа и тело
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "Internal Server Error", w.Body.String())
}

func TestRegisterGetMetricValueHandler_MissingMetricName(t *testing.T) {
	// Создаем тестовый роутер
	router := gin.Default()
	RegisterGetMetricValueHandler(router, new(MockGetValueService))

	// Отправляем запрос только с типом метрики, но без имени
	req, err := http.NewRequest("GET", "/value/gauge", nil)
	assert.NoError(t, err)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Проверяем код ответа и сообщение
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "Metric name is required", w.Body.String())
}
