package handlers

import (
	"bytes"
	"go-metrics-alerting/internal/domain"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для UpdateMetricService
type MockUpdateMetricService struct {
	mock.Mock
}

// Метод UpdateMetric мока, который будет возвращать nil в случае неудачного обновления
func (m *MockUpdateMetricService) UpdateMetric(metric *domain.Metrics) *domain.Metrics {
	args := m.Called(metric)
	// Возвращаем nil, если этот мок настроен так
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*domain.Metrics)
}

// Тест для успешного обновления метрики с Delta
func TestUpdateMetricsBodyHandler_SuccessWithDelta(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockUpdateMetricService)

	// Подготовка данных
	metric := &domain.Metrics{ID: "1", MType: domain.MType("counter"), Value: new(float64)}
	mockService.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(metric)

	r.POST("/update/", UpdateMetricsBodyHandler(mockService))

	// Тело запроса с Delta
	payload := `{"id":"1", "type":"counter", "delta": 10}`

	// Отправка запроса
	req, _ := http.NewRequest("POST", "/update/", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус и тело ответа
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 OK")
	mockService.AssertExpectations(t)
}

// Тест для успешного обновления метрики с Value
func TestUpdateMetricsBodyHandler_SuccessWithValue(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockUpdateMetricService)

	// Подготовка данных
	metric := &domain.Metrics{ID: "1", MType: domain.MType("counter"), Value: new(float64)}
	mockService.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(metric)

	r.POST("/update/", UpdateMetricsBodyHandler(mockService))

	// Тело запроса с Value
	payload := `{"id":"1", "type":"gauge", "value": 10.5}`

	// Отправка запроса
	req, _ := http.NewRequest("POST", "/update/", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус и тело ответа
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 OK")
	mockService.AssertExpectations(t)
}

// Тест для невалидного JSON
func TestUpdateMetricsBodyHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockUpdateMetricService)

	r.POST("/update/", UpdateMetricsBodyHandler(mockService))

	// Тело запроса с ошибкой (некорректный JSON)
	payload := `{"id":"1", "type":"counter", "value":}` // Некорректный JSON

	// Отправка запроса
	req, _ := http.NewRequest("POST", "/update/", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP status 400 Bad Request")
}

// Тест для невалидных данных (ошибка валидации)
func TestUpdateMetricsBodyHandler_InvalidValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockUpdateMetricService)

	// Подготовка данных
	mockService.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(nil)

	r.POST("/update/", UpdateMetricsBodyHandler(mockService))

	// Тело запроса с ошибкой валидации (неправильный тип)
	payload := `{"id":"1", "type":"invalidType", "value": 10.5}`

	// Отправка запроса
	req, _ := http.NewRequest("POST", "/update/", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP status 400 Bad Request")
}

// Тест для успешного обновления метрики с параметрами пути
func TestUpdateMetricsPathHandler_SuccessWithParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockUpdateMetricService)

	// Подготовка данных
	metric := &domain.Metrics{ID: "1", MType: domain.MType("counter"), Value: new(float64)}
	mockService.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(metric)

	r.POST("/update/:mType/:id/:value", UpdateMetricsPathHandler(mockService))

	// Параметры пути: mType = counter, id = 1, value = 10.5
	req, _ := http.NewRequest("POST", "/update/counter/1/10", nil)

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус и тело ответа
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 OK")
	mockService.AssertExpectations(t)
}

// Тест для невалидных данных в параметрах пути (ошибка валидации)
func TestUpdateMetricsPathHandler_InvalidValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockUpdateMetricService)

	// Подготовка данных
	mockService.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(nil)

	r.POST("/update/:mType/:id/:value", UpdateMetricsPathHandler(mockService))

	// Параметры пути с ошибкой: неправильный тип метрики
	req, _ := http.NewRequest("POST", "/update/invalidType/1/10.5", nil)

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа (ошибка валидации)
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP status 400 Bad Request")
}

func TestUpdateMetricsPathHandler_UpdateMetricFailed(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockUpdateMetricService)

	// Настройка ожидания для вызова UpdateMetric
	mockService.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(nil).Once()

	// Регистрация обработчика
	r.POST("/update/:mType/:id/:value", UpdateMetricsPathHandler(mockService))

	// Параметры пути: mType = counter, id = 1, value = 10 (правильный тип для counter)
	req, _ := http.NewRequest("POST", "/update/counter/1/10", nil)

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем, что в ответе содержится сообщение о неудачном обновлении метрики
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP status 400 Bad Request")
	assert.Contains(t, w.Body.String(), "Metric is not updated", "Expected 'Metric is not updated' in response body")

	// Проверяем, что метод UpdateMetric был вызван
	mockService.AssertExpectations(t)
}

// MockGetMetricService - мок для интерфейса GetMetricService
type MockGetMetricService struct {
	mock.Mock
}

// GetMetric - имитация метода получения метрики
func (m *MockGetMetricService) GetMetric(id string, mtype domain.MType) *domain.Metrics {
	args := m.Called(id, mtype)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*domain.Metrics)
}

func TestGetMetricValueBodyHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockGetMetricService)

	// Подготовка данных
	metric := &domain.Metrics{ID: "1", MType: domain.MType("gauge"), Value: new(float64)}
	mockService.On("GetMetric", "1", domain.MType("gauge")).Return(metric)

	// Регистрация обработчика
	r.POST("/get-metric", GetMetricValueBodyHandler(mockService))

	// Тело запроса
	payload := `{"id":"1", "type":"gauge"}`

	// Отправка запроса
	req, _ := http.NewRequest("POST", "/get-metric", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус и тело ответа
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 OK")
	assert.Contains(t, w.Body.String(), "1", "Expected metric ID in the response")
	mockService.AssertExpectations(t)
}

// Тест для невалидного JSON
func TestGetMetricValueBodyHandler_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockGetMetricService)

	// Регистрация обработчика
	r.POST("/get-metric", GetMetricValueBodyHandler(mockService))

	// Некорректный JSON
	payload := `{"id":"1", "type":}` // Некорректный JSON

	// Отправка запроса
	req, _ := http.NewRequest("POST", "/get-metric", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP status 400 Bad Request")
	assert.Contains(t, w.Body.String(), "invalid request", "Expected error message 'invalid request'")
}

// Тест для ошибки валидации данных
func TestGetMetricValueBodyHandler_InvalidValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockGetMetricService)

	// Регистрация обработчика
	r.POST("/get-metric", GetMetricValueBodyHandler(mockService))

	// Тело запроса с ошибкой валидации (неправильный тип метрики)
	payload := `{"id":"1", "type":"invalidType"}`

	// Отправка запроса
	req, _ := http.NewRequest("POST", "/get-metric", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusBadRequest, w.Code, "Expected HTTP status 400 Bad Request")
	assert.Contains(t, w.Body.String(), "error", "Expected error message for invalid validation")
}

// Тест для несуществующей метрики
func TestGetMetricValueBodyHandler_MetricNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockGetMetricService)

	// Настройка мока, чтобы возвращать nil, если метрика не найдена
	mockService.On("GetMetric", "1", domain.MType("counter")).Return(nil)

	// Регистрация обработчика
	r.POST("/get-metric", GetMetricValueBodyHandler(mockService))

	// Тело запроса
	payload := `{"id":"1", "type":"counter"}`

	// Отправка запроса
	req, _ := http.NewRequest("POST", "/get-metric", bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP status 404 Not Found")
	assert.Contains(t, w.Body.String(), "metric not found", "Expected error message 'metric not found'")

	// Проверяем, что метод GetMetric был вызван с ожидаемыми параметрами
	mockService.AssertExpectations(t)
}

// Тест для успешного получения метрики
func TestGetMetricValuePathHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockGetMetricService)

	var v float64 = 10.5
	// Создаем фиктивную метрику
	metric := &domain.Metrics{
		ID:    "1",
		MType: domain.MType("counter"),
		Value: &v,
	}

	// Настройка мока: возвращаем метрику для заданного ID и типа
	mockService.On("GetMetric", "1", domain.MType("counter")).Return(metric)

	// Регистрация обработчика
	r.GET("/get-metric/:mType/:id", GetMetricValuePathHandler(mockService))

	// Отправка запроса
	req, _ := http.NewRequest("GET", "/get-metric/counter/1", nil)

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 OK")
	// Проверяем, что в ответе правильное значение метрики
	assert.Contains(t, w.Body.String(), "10.500000", "Expected response body to contain the metric value")
	mockService.AssertExpectations(t)
}

// Тест для случая, когда метрика не найдена
func TestGetMetricValuePathHandler_MetricNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockGetMetricService)

	// Настройка мока: возвращаем nil, если метрика не найдена
	mockService.On("GetMetric", "1", domain.MType("counter")).Return(nil)

	// Регистрация обработчика
	r.GET("/get-metric/:mType/:id", GetMetricValuePathHandler(mockService))

	// Отправка запроса
	req, _ := http.NewRequest("GET", "/get-metric/counter/1", nil)

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP status 404 Not Found")
	// Проверяем, что в ответе содержится сообщение "Metric not found"
	assert.Contains(t, w.Body.String(), "Metric not found", "Expected response body to contain 'Metric not found'")
	mockService.AssertExpectations(t)
}

func TestGetMetricValuePathHandler_InvalidMType(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockGetMetricService)

	// Настройка мока: возвращаем nil для всех запросов с типом "unknown"
	mockService.On("GetMetric", "1", domain.MType("unknown")).Return(nil)

	// Регистрация обработчика
	r.GET("/get-metric/:mType/:id", GetMetricValuePathHandler(mockService))

	// Отправка запроса с неверным типом метрики (например, "unknown" вместо "counter" или "gauge")
	req, _ := http.NewRequest("GET", "/get-metric/unknown/1", nil)

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа (должен быть 404)
	assert.Equal(t, http.StatusNotFound, w.Code, "Expected HTTP status 404 Not Found")
	// Проверяем, что в ответе содержится сообщение "Metric not found"
	assert.Contains(t, w.Body.String(), "Metric not found", "Expected response body to contain 'Metric not found'")
	mockService.AssertExpectations(t)
}

// MockGetAllMetricsService - мок для GetAllMetricsService
type MockGetAllMetricsService struct {
	mock.Mock
}

// GetAllMetrics имитирует получение всех метрик
func (m *MockGetAllMetricsService) GetAllMetrics() []*domain.Metrics {
	args := m.Called()
	return args.Get(0).([]*domain.Metrics) // Возвращаем метки как срез указателей на domain.Metrics
}

// Тест для успешного получения всех метрик и их отображения в HTML
func TestGetAllMetricValuesHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Создаем новый роутер
	r := gin.Default()

	// Мокаем сервис
	mockService := new(MockGetAllMetricsService)

	// Подготовка данных
	metrics := []*domain.Metrics{
		{ID: "1", MType: "counter", Value: new(float64)},
		{ID: "2", MType: "gauge", Value: new(float64)},
	}

	// Задаем поведение мока
	mockService.On("GetAllMetrics").Return(metrics)

	// Регистрация обработчика
	r.GET("/get-all-metrics", GetAllMetricValuesHandler(mockService))

	// Отправка запроса
	req, _ := http.NewRequest("GET", "/get-all-metrics", nil)

	// Записываем ответ
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Проверяем статус ответа (должен быть 200 OK)
	assert.Equal(t, http.StatusOK, w.Code, "Expected HTTP status 200 OK")

	// Проверяем, что метрики были отображены в ответе HTML
	// Проверка, что ID "1" присутствует в ответе
	assert.True(t, strings.Contains(w.Body.String(), "1: "), "Expected metric ID '1' in the response")
	// Проверка, что ID "2" присутствует в ответе
	assert.True(t, strings.Contains(w.Body.String(), "2: "), "Expected metric ID '2' in the response")

	// Проверяем, что нет сообщения 'No metrics available'
	assert.False(t, strings.Contains(w.Body.String(), "No metrics available"), "Expected metrics to be displayed, not 'No metrics available'")

	// Проверяем, что метод GetAllMetrics был вызван
	mockService.AssertExpectations(t)
}
