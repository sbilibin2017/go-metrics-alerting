package handlers

import (
	"go-metrics-alerting/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUpdateMetricService мокируем сервис обновления метрик
type MockUpdateMetricPathService struct {
	mock.Mock
}

func (m *MockUpdateMetricPathService) UpdateMetric(metric *domain.Metrics) (*domain.Metrics, error) {
	args := m.Called(metric)
	// Если ошибка, возвращаем nil для метрики и ошибку
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	// Иначе возвращаем метрику
	return args.Get(0).(*domain.Metrics), nil
}

// Тест для успешного обновления метрики
func TestUpdateMetricPathHandler_Success(t *testing.T) {
	// Создаем мок-сервис
	mockService := new(MockUpdateMetricPathService)
	metric := &domain.Metrics{
		ID:    "test",
		MType: "counter", // Пример типа метрики
		Value: float64Ptr(25.0),
	}
	mockService.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(metric, nil)

	// Создаем новый роутер
	r := chi.NewRouter()
	r.Post("/update/{type}/{id}/{value}", UpdateMetricPathHandler(mockService))

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Создаем запрос для успешного обновления
	req, err := http.NewRequest("POST", ts.URL+"/update/counter/test/25", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус код ответа
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Проверяем, что метод UpdateMetric был вызван
	mockService.AssertExpectations(t)
}

// Тест для запроса с пустым ID
func TestUpdateMetricPathHandler_MissingID(t *testing.T) {
	// Создаем мок-сервис
	mockService := new(MockUpdateMetricPathService)

	// Создаем новый роутер
	r := chi.NewRouter()
	r.Post("/update/{type}/{id}/{value}", UpdateMetricPathHandler(mockService))

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Создаем запрос с пустым ID
	req, err := http.NewRequest("POST", ts.URL+"/update/counter//25", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус код ответа
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Тест для запроса с пустым типом
func TestUpdateMetricPathHandler_MissingType(t *testing.T) {
	// Создаем мок-сервис
	mockService := new(MockUpdateMetricPathService)

	// Создаем новый роутер
	r := chi.NewRouter()
	r.Post("/update/{type}/{id}/{value}", UpdateMetricPathHandler(mockService))

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Создаем запрос с пустым типом
	req, err := http.NewRequest("POST", ts.URL+"/update//test/25", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус код ответа
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Тест для запроса с пустым значением
func TestUpdateMetricPathHandler_MissingValue(t *testing.T) {
	// Создаем мок-сервис
	mockService := new(MockUpdateMetricPathService)

	// Создаем новый роутер
	r := chi.NewRouter()
	r.Post("/update/{type}/{id}/{value}", UpdateMetricPathHandler(mockService))

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Создаем запрос с пустым значением
	req, err := http.NewRequest("POST", ts.URL+"/update/counter/test", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус код ответа
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

// Тест для запроса с некорректным типом метрики
func TestUpdateMetricPathHandler_InvalidMetricType(t *testing.T) {
	// Создаем мок-сервис
	mockService := new(MockUpdateMetricPathService)

	// Создаем новый роутер
	r := chi.NewRouter()
	r.Post("/update/{type}/{id}/{value}", UpdateMetricPathHandler(mockService))

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Создаем запрос с некорректным типом метрики
	req, err := http.NewRequest("POST", ts.URL+"/update/invalid_type/test/25", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус код ответа
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

// Тест для ошибки при обновлении метрики
func TestUpdateMetricPathHandler_Error(t *testing.T) {
	// Создаем мок-сервис
	mockService := new(MockUpdateMetricPathService)
	mockService.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(nil, assert.AnError)

	// Создаем новый роутер
	r := chi.NewRouter()
	r.Post("/update/{type}/{id}/{value}", UpdateMetricPathHandler(mockService))

	// Создаем тестовый сервер
	ts := httptest.NewServer(r)
	defer ts.Close()

	// Создаем запрос для обновления
	req, err := http.NewRequest("POST", ts.URL+"/update/counter/test/25", nil)
	assert.NoError(t, err)

	// Отправляем запрос
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Проверяем статус код ответа
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	// Проверяем, что метод UpdateMetric был вызван
	mockService.AssertExpectations(t)
}

// Утилита для создания указателя на float64
func float64Ptr(v float64) *float64 {
	return &v
}
