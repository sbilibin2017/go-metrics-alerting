package services

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/domain"

	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// Мок для коллекционера метрик
type MockMetricsCollector struct {
	mock.Mock
}

func (m *MockMetricsCollector) Collect() []*domain.Metrics {
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

func TestMetricAgent_Run(t *testing.T) {
	// Настройка мока для коллекционера
	mockCollector := new(MockMetricsCollector)
	mockFacade := new(MockMetricFacade)

	// Пример метрики для коллекции
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			MType: domain.Counter,
			Value: float64Ptr(5),
		},
	}
	mockCollector.On("Collect").Return(metrics)

	// Настройка ожидания на вызов метода UpdateMetric
	mockFacade.On("UpdateMetric", mock.AnythingOfType("*domain.Metrics")).Return(nil)

	// Настройка конфигурации
	config := &configs.AgentConfig{
		Address:        "localhost:8080",
		PollInterval:   100 * time.Millisecond,
		ReportInterval: 200 * time.Millisecond,
	}

	// Настройка логгера
	logger, _ := zap.NewProduction()

	// Создание агента
	agent := NewMetricAgentService(config, []MetricsCollector{mockCollector}, mockFacade, logger)

	// Канал для имитации сигнала завершения
	signalCh := make(chan os.Signal, 1)

	// Запуск агента в отдельной горутине
	go agent.Run(signalCh)

	// Пождать несколько секунд, чтобы агент успел собрать и отправить метрики
	time.Sleep(1 * time.Second)

	// Проверка вызовов
	mockCollector.AssertExpectations(t)
	mockFacade.AssertExpectations(t)

	// Проверка вызова collect и send
	mockCollector.AssertCalled(t, "Collect")
	mockFacade.AssertCalled(t, "UpdateMetric", mock.AnythingOfType("*domain.Metrics"))
}

func TestMetricAgent_Run_ShutdownSignal(t *testing.T) {
	// Настройка мока для коллекционера
	mockCollector := new(MockMetricsCollector)
	mockFacade := new(MockMetricFacade)

	// Пример метрики для коллекции
	metrics := []*domain.Metrics{
		{
			ID:    "metric1",
			MType: domain.Counter,
			Value: float64Ptr(5),
		},
	}
	mockCollector.On("Collect").Return(metrics)

	// Настройка конфигурации
	config := &configs.AgentConfig{
		Address:        "localhost:8080",
		PollInterval:   100 * time.Millisecond,
		ReportInterval: 200 * time.Millisecond,
	}

	// Использование zaptest.NewLogger для проверки логов
	logger := zaptest.NewLogger(t)

	// Создание агента
	agent := NewMetricAgentService(config, []MetricsCollector{mockCollector}, mockFacade, logger)

	// Канал для имитации сигнала завершения
	signalCh := make(chan os.Signal, 1)

	// Запуск агента в горутине
	go agent.Run(signalCh)

	// Отправляем сигнал о завершении
	signalCh <- os.Interrupt

	// Пождать, пока агент обработает сигнал и завершится
	time.Sleep(500 * time.Millisecond)

	// Проверка, что логгер вызвал сообщение о завершении

	mockFacade.AssertExpectations(t)
}

// Утилита для создания указателя на float64
func float64Ptr(f float64) *float64 {
	return &f
}
