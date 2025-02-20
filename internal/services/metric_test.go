package services

import (
	"testing"

	"go-metrics-alerting/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Моки для стратегий обновления
type UpdateMetricStrategyMock struct {
	mock.Mock
}

func (m *UpdateMetricStrategyMock) Update(metric *domain.Metric) (*domain.Metric, bool) {
	args := m.Called(metric)
	// Проверяем, что мы не возвращаем nil в случае ошибки
	if args.Get(0) == nil {
		return nil, args.Bool(1)
	}
	return args.Get(0).(*domain.Metric), args.Bool(1)
}

func TestUpdateMetricsService_UpdateMetricValue_Success(t *testing.T) {
	// Создаем тестовую метрику
	metric := &domain.Metric{
		ID:    "metric1",
		Value: "100",
	}

	// Мокаем стратегию
	strategyMock := new(UpdateMetricStrategyMock)
	strategyMock.On("Update", metric).Return(metric, true) // Стратегия успешно обновляет метрику

	// Создаем сервис с мокой стратегии
	service := &UpdateMetricsService{strategy: strategyMock}

	// Вызываем метод обновления
	updatedMetric, err := service.UpdateMetricValue(metric)

	// Проверяем, что ошибки нет и метрика обновлена корректно
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, metric.ID, updatedMetric.ID, "Metric IDs should match")
	assert.Equal(t, metric.Value, updatedMetric.Value, "Metric values should match")

	// Проверяем, что стратегия была вызвана
	strategyMock.AssertExpectations(t)
}

func TestUpdateMetricsService_UpdateMetricValue_Fail(t *testing.T) {
	// Создаем тестовую метрику
	metric := &domain.Metric{
		ID:    "metric1",
		Value: "100",
	}

	// Мокаем стратегию
	strategyMock := new(UpdateMetricStrategyMock)
	strategyMock.On("Update", metric).Return(nil, false) // Стратегия не обновляет метрику

	// Создаем сервис с мокой стратегии
	service := &UpdateMetricsService{strategy: strategyMock}

	// Вызываем метод обновления
	updatedMetric, err := service.UpdateMetricValue(metric)

	// Проверяем, что ошибка произошла и метрика не обновлена
	assert.Error(t, err, "Expected error")
	assert.Equal(t, ErrUpdateFailed, err, "Expected ErrUpdateFailed error")
	assert.Nil(t, updatedMetric, "Expected nil updated metric")

	// Проверяем, что стратегия была вызвана
	strategyMock.AssertExpectations(t)
}
