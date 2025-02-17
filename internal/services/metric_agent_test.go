package services

import (
	"errors"
	"testing"
	"time"

	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/types"

	"github.com/stretchr/testify/mock"
)

// Мок фасада метрик
type MockMetricFacade struct {
	mock.Mock
}

func (m *MockMetricFacade) UpdateMetric(metric types.MetricsRequest) error {
	args := m.Called(metric)
	return args.Error(0)
}

// Мок сборщика метрик
type MockMetricCollector struct {
	mock.Mock
}

func (m *MockMetricCollector) Collect() []types.MetricsRequest {
	args := m.Called()
	return args.Get(0).([]types.MetricsRequest)
}

// Вспомогательные функции для удобства
func int64Ptr(i int64) *int64 {
	return &i
}

func float64Ptr(f float64) *float64 {
	return &f
}

func TestMetricAgentService_SuccessfulCollectionAndReporting(t *testing.T) {
	// Мок фасада
	mockFacade := new(MockMetricFacade)
	mockCollector := new(MockMetricCollector)

	mockMetrics := []types.MetricsRequest{
		{
			ID:    "1",
			MType: "gauge",
			Value: float64Ptr(10.5),
		},
		{
			ID:    "2",
			MType: "counter",
			Delta: int64Ptr(5),
		},
	}

	// Настроим коллектора на возврат тестовых метрик
	mockCollector.On("Collect").Return(mockMetrics)

	// Настроим фасад для успешного обновления метрик
	for _, metric := range mockMetrics {
		mockFacade.On("UpdateMetric", metric).Return(nil)
	}

	service := NewMetricAgentService(
		&configs.AgentConfig{
			PollInterval:   1,
			ReportInterval: 1,
		},
		mockFacade,
		[]MetricCollector{mockCollector},
	)

	// Стартуем сервис в отдельной горутине
	go service.Start()

	// Ждем немного времени, чтобы процесс сбора и отправки метрик успел завершиться
	time.Sleep(2 * time.Second)

	// Проверяем, что фасад был вызван с нужными метками
	for _, metric := range mockMetrics {
		mockFacade.AssertCalled(t, "UpdateMetric", metric)
	}
}

func TestMetricAgentService_ErrorInUpdatingMetric(t *testing.T) {
	// Мок фасада
	mockFacade := new(MockMetricFacade)
	mockCollector := new(MockMetricCollector)

	mockMetrics := []types.MetricsRequest{
		{
			ID:    "1",
			MType: "gauge",
			Value: float64Ptr(10.5),
		},
	}

	// Настроим коллектора на возврат тестовых метрик
	mockCollector.On("Collect").Return(mockMetrics)

	// Настроим фасад для возврата ошибки при обновлении метрики
	mockFacade.On("UpdateMetric", mockMetrics[0]).Return(errors.New("some error"))

	service := NewMetricAgentService(
		&configs.AgentConfig{
			PollInterval:   1,
			ReportInterval: 1,
		},
		mockFacade,
		[]MetricCollector{mockCollector},
	)

	// Стартуем сервис в отдельной горутине
	go service.Start()

	// Ждем немного времени, чтобы процесс сбора и отправки метрик успел завершиться
	time.Sleep(2 * time.Second)

	// Проверяем, что фасад был вызван с нужными метками
	mockFacade.AssertCalled(t, "UpdateMetric", mockMetrics[0])

	// Проверяем, что ошибка была обработана
	mockFacade.AssertExpectations(t)
}
