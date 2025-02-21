package services

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/domain"

	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

// Мок для стратегии коллекционирования метрик
type MockMetricsCollectorStrategy struct {
	mock.Mock
}

func (m *MockMetricsCollectorStrategy) Collect() []*domain.Metrics {
	args := m.Called()
	return args.Get(0).([]*domain.Metrics)
}

// Мок для фасада метрик
type MockMetricFacade struct {
	mock.Mock
}

func (m *MockMetricFacade) UpdateMetric(metric *domain.Metrics) {
	m.Called(metric)
}

func TestMetricAgentService_Run(t *testing.T) {
	// Создаем моки
	mockGaugeCollector := new(MockMetricsCollectorStrategy)
	mockCounterCollector := new(MockMetricsCollectorStrategy)
	mockFacade := new(MockMetricFacade)

	// Пример метрики
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			MType: domain.Gauge,
			Value: float64Ptr(5),
		},
	}

	// Настройка моков
	mockGaugeCollector.On("Collect").Return(metrics)
	mockCounterCollector.On("Collect").Return([]*domain.Metrics{})
	mockFacade.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(nil)

	// Конфигурация агента
	config := &configs.AgentConfig{
		Address:        "localhost:8080",
		PollInterval:   100 * time.Millisecond,
		ReportInterval: 200 * time.Millisecond,
	}

	// Используем zaptest.Logger
	logger := zaptest.NewLogger(t)

	// Создаем мапу стратегий
	collectorStrategies := map[domain.MType]MetricsCollectorStrategy{
		domain.Gauge:   mockGaugeCollector,
		domain.Counter: mockCounterCollector,
	}

	// Создаем агент
	agent := NewMetricAgentService(config, collectorStrategies, mockFacade, logger)

	// Канал для сигнала завершения
	signalCh := make(chan os.Signal, 1)

	// Запуск агента в горутине
	go agent.Run(signalCh)

	// Ждем, чтобы агент успел собрать и отправить метрики
	time.Sleep(1 * time.Second)

	// Проверка вызовов
	mockGaugeCollector.AssertExpectations(t)
	mockCounterCollector.AssertExpectations(t)
	mockFacade.AssertExpectations(t)

	// Проверяем, что метод Collect был вызван
	mockGaugeCollector.AssertCalled(t, "Collect")
	mockCounterCollector.AssertCalled(t, "Collect")
	mockFacade.AssertCalled(t, "UpdateMetric", mock.AnythingOfType("*domain.Metrics"))
}

func TestMetricAgentService_Run_ShutdownSignal(t *testing.T) {
	// Создаем моки
	mockGaugeCollector := new(MockMetricsCollectorStrategy)
	mockCounterCollector := new(MockMetricsCollectorStrategy)
	mockFacade := new(MockMetricFacade)

	// Пример метрики
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			MType: domain.Gauge,
			Value: float64Ptr(5),
		},
	}

	// Настройка моков
	mockGaugeCollector.On("Collect").Return(metrics)
	mockCounterCollector.On("Collect").Return([]*domain.Metrics{})

	// Конфигурация агента
	config := &configs.AgentConfig{
		Address:        "localhost:8080",
		PollInterval:   100 * time.Millisecond,
		ReportInterval: 200 * time.Millisecond,
	}

	// Используем zaptest.Logger
	logger := zaptest.NewLogger(t)

	// Создаем мапу стратегий
	collectorStrategies := map[domain.MType]MetricsCollectorStrategy{
		domain.Gauge:   mockGaugeCollector,
		domain.Counter: mockCounterCollector,
	}

	// Создаем агент
	agent := NewMetricAgentService(config, collectorStrategies, mockFacade, logger)

	// Канал для сигнала завершения
	signalCh := make(chan os.Signal, 1)

	// Запуск агента в горутине
	go agent.Run(signalCh)

	// Отправляем сигнал завершения
	signalCh <- os.Interrupt

	// Ждем, пока агент завершится
	time.Sleep(500 * time.Millisecond)

	// Проверяем, что мокированный логгер обработал завершение корректно
	mockFacade.AssertExpectations(t)
}

// Утилита для создания указателя на float64
func float64Ptr(f float64) *float64 {
	return &f
}
