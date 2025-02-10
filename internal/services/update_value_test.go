package services

import (
	"context"
	e "errors"
	"go-metrics-alerting/internal/errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockMetricRepository — мок-реализация интерфейса MetricRepository
type MockMetricRepository struct {
	mock.Mock
}

func (m *MockMetricRepository) Save(ctx context.Context, metricType, metricName, value string) error {
	args := m.Called(ctx, metricType, metricName, value)
	return args.Error(0)
}

func (m *MockMetricRepository) Get(ctx context.Context, metricType, metricName string) (string, error) {
	args := m.Called(ctx, metricType, metricName)
	return args.String(0), args.Error(1)
}

// Тест на ошибку при сохранении метрики
func TestUpdateMetricValueService_SaveError(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  types.Counter,
		Name:  "metric1",
		Value: "5",
	}

	// Мокаем успешное получение текущего значения метрики
	mockRepo.On("Get", mock.Anything, string(types.Counter), "metric1").Return("10", nil)
	// Мокаем ошибку при сохранении
	mockRepo.On("Save", mock.Anything, string(types.Counter), "metric1", "15").Return(e.New("save error"))

	err := service.UpdateMetricValue(context.Background(), req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)
	assert.Equal(t, errors.ErrValueNotSaved.Error(), apiErr.Message)
	mockRepo.AssertExpectations(t)
}

// Позитивный тест на успешное обновление значения метрики
func TestUpdateMetricValueService_Success(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{MetricRepository: mockRepo}

	req := &types.UpdateMetricValueRequest{
		Type:  types.Counter,
		Name:  "metric1",
		Value: "5",
	}

	// Мокаем успешное получение текущего значения метрики
	mockRepo.On("Get", mock.Anything, string(types.Counter), "metric1").Return("10", nil)
	// Мокаем успешное сохранение нового значения
	mockRepo.On("Save", mock.Anything, string(types.Counter), "metric1", "15").Return(nil)

	err := service.UpdateMetricValue(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Тест на обработку Gauge метрики
func TestUpdateMetricValueService_GaugeProcessing(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{
		MetricRepository: mockRepo,
	}

	req := &types.UpdateMetricValueRequest{
		Type:  types.Gauge,
		Name:  "metric2",
		Value: "5.5",
	}

	// Мокаем успешное получение текущего значения метрики
	mockRepo.On("Get", mock.Anything, string(types.Gauge), "metric2").Return("3.3", nil)

	// Выполняем преобразование строки в число
	newVal, _ := strconv.ParseFloat(req.Value, 64)
	value := strconv.FormatFloat(newVal, 'f', -1, 64)

	// Мокаем успешное сохранение нового значения
	mockRepo.On("Save", mock.Anything, string(types.Gauge), "metric2", value).Return(nil)

	err := service.UpdateMetricValue(context.Background(), req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Тест на ошибку при получении значения метрики с установкой значения по умолчанию
func TestUpdateMetricValueService_GetErrorHandling(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := &UpdateMetricValueService{
		MetricRepository: mockRepo,
	}

	req := &types.UpdateMetricValueRequest{
		Type:  types.Counter,
		Name:  "metric1",
		Value: "5",
	}

	// Мокаем ошибку при получении текущего значения метрики
	mockRepo.On("Get", mock.Anything, string(types.Counter), "metric1").Return("", e.New("get error"))

	// Мокаем успешное сохранение нового значения
	mockRepo.On("Save", mock.Anything, string(types.Counter), "metric1", "5").Return(nil)

	// Выполняем обновление метрики
	err := service.UpdateMetricValue(context.Background(), req)

	// Проверяем, что ошибка отсутствует, а значение установлено в "0"
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
