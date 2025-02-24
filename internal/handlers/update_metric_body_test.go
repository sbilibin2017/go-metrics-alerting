package handlers

import (
	"bytes"
	"encoding/json"
	"go-metrics-alerting/internal/domain"
	"go-metrics-alerting/internal/types"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mailru/easyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUpdateMetricService мокируем сервис обновления метрик
type MockUpdateMetricBodyService struct {
	mock.Mock
}

func (m *MockUpdateMetricBodyService) UpdateMetric(metric *domain.Metrics) (*domain.Metrics, error) {
	args := m.Called(metric)
	// Если ошибка, возвращаем nil для метрики и ошибку
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	// Иначе возвращаем метрику
	return args.Get(0).(*domain.Metrics), nil
}

func TestUpdateMetricBodyHandler_Success(t *testing.T) {
	// Создаем мок сервиса
	mockService := new(MockUpdateMetricBodyService)

	// Подготавливаем тестовые данные
	reqBody := types.UpdateMetricBodyRequest{
		ID:    "123",
		MType: string(domain.Counter),
		Delta: new(int64),
	}

	*reqBody.Delta = 10

	// Мокаем поведение метода UpdateMetric
	mockService.On("UpdateMetric", mock.Anything).Return(&domain.Metrics{
		ID:    reqBody.ID,
		MType: domain.Counter,
		Delta: reqBody.Delta,
		Value: nil,
	}, nil)

	// Создаем новый запрос с тестовыми данными
	reqBodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/update-metric", bytes.NewReader(reqBodyJSON))

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Создаем обработчик и передаем мок-сервис
	handler := UpdateMetricBodyHandler(mockService)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusOK, rr.Code)

	// Проверяем правильность ответа
	var response types.UpdateMetricBodyResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, reqBody.ID, response.ID)
	assert.Equal(t, reqBody.MType, response.MType)
	assert.Equal(t, *reqBody.Delta, *response.Delta)

	// Проверяем, что метод UpdateMetric был вызван
	mockService.AssertExpectations(t)
}

func TestUpdateMetricBodyHandler_InvalidRequest(t *testing.T) {
	// Создаем мок сервиса
	mockService := new(MockUpdateMetricBodyService)

	// Создаем некорректный запрос
	req := httptest.NewRequest(http.MethodPost, "/update-metric", nil)
	rr := httptest.NewRecorder()

	// Создаем обработчик и передаем мок-сервис
	handler := UpdateMetricBodyHandler(mockService)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateMetricBodyHandler_ValidationError(t *testing.T) {
	// Создаем мок сервиса
	mockService := new(MockUpdateMetricBodyService)

	// Создаем запрос с ошибкой валидации (например, пустой ID)
	reqBody := types.UpdateMetricBodyRequest{
		ID:    "",
		MType: string(domain.Counter),
		Delta: new(int64),
	}

	*reqBody.Delta = 10

	reqBodyJSON, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/update-metric", bytes.NewReader(reqBodyJSON))
	rr := httptest.NewRecorder()

	// Создаем обработчик и передаем мок-сервис
	handler := UpdateMetricBodyHandler(mockService)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем статус код
	assert.Equal(t, http.StatusNotFound, rr.Code)
}

// Тест для ошибки в методе UpdateMetric
func TestUpdateMetricBodyHandler_UpdateMetricError(t *testing.T) {
	// Создаем мок сервиса
	mockService := new(MockUpdateMetricBodyService)

	// Создаем тестовые данные
	reqBody := types.UpdateMetricBodyRequest{
		ID:    "123",
		MType: string(domain.Counter),
		Delta: new(int64),
	}

	*reqBody.Delta = 10

	// Мокаем ошибку, которую возвращает метод UpdateMetric
	mockService.On("UpdateMetric", mock.Anything).Return(nil, assert.AnError)

	// Создаем новый запрос с тестовыми данными
	reqBodyJSON, _ := easyjson.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/update-metric", bytes.NewReader(reqBodyJSON))

	// Создаем ResponseRecorder для записи ответа
	rr := httptest.NewRecorder()

	// Создаем обработчик и передаем мок-сервис
	handler := UpdateMetricBodyHandler(mockService)

	// Выполняем запрос
	handler.ServeHTTP(rr, req)

	// Проверяем, что возвращен статус 500 (Internal Server Error)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Проверяем, что в теле ответа содержится текст ошибки
	expectedErrorMessage := assert.AnError.Error()
	assert.Contains(t, rr.Body.String(), expectedErrorMessage)

	// Проверяем, что метод UpdateMetric был вызван
	mockService.AssertExpectations(t)
}
