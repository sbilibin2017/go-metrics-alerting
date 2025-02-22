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
	gaugeMetrics := []*domain.Metrics{
		{
			ID:    "gaugeMetric1",
			MType: domain.Gauge,
			Value: float64Ptr(5),
		},
	}
	counterMetrics := []*domain.Metrics{
		{
			ID:    "counterMetric1",
			MType: domain.Counter,
			Value: float64Ptr(10),
		},
	}

	// Настройка моков
	mockGaugeCollector.On("Collect").Return(gaugeMetrics)
	mockCounterCollector.On("Collect").Return(counterMetrics)
	mockFacade.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(nil)

	// Конфигурация агента
	config := &configs.AgentConfig{
		Address:        "localhost:8080",
		PollInterval:   100 * time.Millisecond,
		ReportInterval: 200 * time.Millisecond,
	}

	// Используем zaptest.Logger
	logger := zaptest.NewLogger(t)

	// Создаем агент с двумя стратегиями коллекции
	agent := NewMetricAgentService(config, mockCounterCollector, mockGaugeCollector, mockFacade, logger)

	// Канал для сигнала завершения
	signalCh := make(chan os.Signal, 1)

	// Запуск агента в горутине
	go agent.Run(signalCh)

	// Ждем, чтобы агент успел собрать и отправить метрики
	time.Sleep(1 * time.Second)

	// Проверка вызовов для первого и второго коллекционера
	mockGaugeCollector.AssertExpectations(t)
	mockCounterCollector.AssertExpectations(t)
	mockFacade.AssertExpectations(t)

	// Проверяем, что метод Collect был вызван для первого и второго коллекционера
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
	gaugeMetrics := []*domain.Metrics{
		{
			ID:    "gaugeMetric1",
			MType: domain.Gauge,
			Value: float64Ptr(5),
		},
	}
	counterMetrics := []*domain.Metrics{
		{
			ID:    "counterMetric1",
			MType: domain.Counter,
			Value: float64Ptr(10),
		},
	}

	// Настройка моков
	mockGaugeCollector.On("Collect").Return(gaugeMetrics)
	mockCounterCollector.On("Collect").Return(counterMetrics)

	// Конфигурация агента
	config := &configs.AgentConfig{
		Address:        "localhost:8080",
		PollInterval:   100 * time.Millisecond,
		ReportInterval: 200 * time.Millisecond,
	}

	// Используем zaptest.Logger
	logger := zaptest.NewLogger(t)

	// Создаем агент с двумя стратегиями коллекции
	agent := NewMetricAgentService(config, mockCounterCollector, mockGaugeCollector, mockFacade, logger)

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
