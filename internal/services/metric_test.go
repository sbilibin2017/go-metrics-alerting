package services

import (
	"context"
	"errors"
	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/pkg/apierror"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func (m *MockMetricRepository) GetAll(ctx context.Context) [][]string {
	args := m.Called(ctx)
	return args.Get(0).([][]string)
}

func TestUpdateMetric_MissingMetricName(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	req := &types.UpdateMetricValueRequest{Name: "", Type: types.Counter, Value: "5"}
	err := service.UpdateMetric(ctx, req)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
}

func TestUpdateMetric_MissingMetricType(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: "", Value: "5"}
	err := service.UpdateMetric(ctx, req)
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
}

func TestUpdateMetric_GetMetricError_DefaultValueUsed(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "counter", "test_metric").Return("", errors.New("not found"))
	mockRepo.On("Save", ctx, "counter", "test_metric", "5").Return(nil)

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Counter, Value: "5"}
	err := service.UpdateMetric(ctx, req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateMetric_InvalidCounterValue(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "counter", "test_metric").Return("", nil)

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Counter, Value: "invalid"}
	err := service.UpdateMetric(ctx, req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMetric_InvalidGaugeValue(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "gauge", "test_metric").Return("", nil)

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Gauge, Value: "invalid"}
	err := service.UpdateMetric(ctx, req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMetric_SuccessfulCounterUpdate(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "counter", "test_metric").Return("10", nil)
	mockRepo.On("Save", ctx, "counter", "test_metric", "15").Return(nil)

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Counter, Value: "5"}
	err := service.UpdateMetric(ctx, req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetMetric_Success(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "counter", "test_metric").Return("10", nil)
	req := &types.GetMetricValueRequest{Name: "test_metric", Type: "counter"}
	value, err := service.GetMetric(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, "10", value)
}

func TestGetMetric_NotFound(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "counter", "test_metric").Return("", errors.New("not found"))
	req := &types.GetMetricValueRequest{Name: "test_metric", Type: "counter"}
	value, err := service.GetMetric(ctx, req)
	assert.Error(t, err)
	assert.Equal(t, "", value)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
}

func TestListMetrics_Success(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("List", ctx).Return([][]string{
		{"counter", "test_metric", "10"},
		{"gauge", "cpu_usage", "99.9"},
	})

	metrics := service.ListMetrics(ctx)
	assert.Len(t, metrics, 2)
	assert.Equal(t, "test_metric", metrics[0].Name)
	assert.Equal(t, "10", metrics[0].Value)
	assert.Equal(t, "cpu_usage", metrics[1].Name)
	assert.Equal(t, "99.9", metrics[1].Value)
}

func TestListMetrics_EmptyList(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("List", ctx).Return([][]string{})

	metrics := service.ListMetrics(ctx)
	assert.Len(t, metrics, 0)
}

func TestUpdateMetric_InvalidMetricType(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: "invalid_type", Value: "10"}
	err := service.UpdateMetric(ctx, req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "invalid metric type", apiErr.Message)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMetric_SaveError(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "counter", "test_metric").Return("10", nil)
	mockRepo.On("Save", ctx, "counter", "test_metric", "15").Return(errors.New("save error"))

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Counter, Value: "5"}
	err := service.UpdateMetric(ctx, req)

	assert.Error(t, err)

	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)
	assert.Equal(t, "value is not saved", apiErr.Message)

	mockRepo.AssertExpectations(t)
}

// Additional Test Cases

func TestUpdateMetric_EmptyMetricType(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: "", Value: "5"}
	err := service.UpdateMetric(ctx, req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusNotFound, apiErr.Code)
	assert.Equal(t, "metric type is required", apiErr.Message)
}

func TestUpdateMetric_InvalidGaugeValueFormat(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	// Настроим мок для метода Get, чтобы он не вызвал панику
	mockRepo.On("Get", ctx, "gauge", "test_metric").Return("", nil) // Возвращаем пустое значение и отсутствие ошибки

	// Теперь тестируем обновление метрики с некорректным значением
	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Gauge, Value: "invalid_value"}
	err := service.UpdateMetric(ctx, req)

	// Проверка на ошибку и правильный HTTP статус
	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "invalid gauge value", apiErr.Message)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMetric_NegativeCounterValue(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "counter", "test_metric").Return("10", nil)
	mockRepo.On("Save", ctx, "counter", "test_metric", "5").Return(nil)

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Counter, Value: "-5"}
	err := service.UpdateMetric(ctx, req)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListMetrics_NilList(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	// Настроим мок для метода List, чтобы он возвращал пустой срез
	mockRepo.On("List", ctx).Return([][]string{})

	// Теперь тестируем получение списка метрик
	metrics := service.ListMetrics(ctx)

	// Проверка, что список пуст
	assert.Len(t, metrics, 0)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMetric_SaveError_DefaultZero(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	mockRepo.On("Get", ctx, "counter", "test_metric").Return("", nil)
	mockRepo.On("Save", ctx, "counter", "test_metric", "0").Return(errors.New("save error"))

	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Counter, Value: "0"}
	err := service.UpdateMetric(ctx, req)

	assert.Error(t, err)
	apiErr, ok := err.(*apierror.APIError)
	assert.True(t, ok)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Code)
	assert.Equal(t, "value is not saved", apiErr.Message)

	mockRepo.AssertExpectations(t)
}

func TestUpdateMetric_SuccessfulGaugeFormat(t *testing.T) {
	mockRepo := new(MockMetricRepository)
	service := MetricService{MetricRepository: mockRepo}
	ctx := context.Background()

	// Настроим мок для метода Get, чтобы он возвращал пустое значение
	mockRepo.On("Get", ctx, "gauge", "test_metric").Return("", nil)

	// Настроим мок для метода Save, чтобы он успешно сохранял метрику
	mockRepo.On("Save", ctx, "gauge", "test_metric", "99.9").Return(nil)

	// Создаем запрос с корректным значением для Gauge
	req := &types.UpdateMetricValueRequest{Name: "test_metric", Type: types.Gauge, Value: "99.9"}

	// Тестируем обновление метрики
	err := service.UpdateMetric(ctx, req)

	// Проверяем, что ошибка отсутствует
	assert.NoError(t, err)

	// Проверяем, что метод Save был вызван с правильными параметрами
	mockRepo.AssertExpectations(t)
}
