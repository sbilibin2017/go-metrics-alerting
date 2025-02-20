package services

import (
	"go-metrics-alerting/internal/configs"
	"go-metrics-alerting/internal/domain"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// Mock для интерфейса MetricsCollector
type MockMetricsCollector struct {
	mock.Mock
}

func (m *MockMetricsCollector) Collect() []*domain.Metric {
	args := m.Called()
	return args.Get(0).([]*domain.Metric)
}

// Mock для интерфейса MetricFacade
type MockMetricFacade struct {
	mock.Mock
}

func (m *MockMetricFacade) UpdateMetric(metric *domain.Metric) {
	m.Called(metric)
}

func TestMetricAgent_Run(t *testing.T) {
	tests := []struct {
		name            string
		mockCollectors  []*MockMetricsCollector
		mockFacade      *MockMetricFacade
		expectedMetrics []*domain.Metric
		collectCalled   bool
		sendCalled      bool
		interruptSignal bool
	}{
		{
			name: "should collect and send metrics",
			mockCollectors: []*MockMetricsCollector{
				&MockMetricsCollector{},
			},
			mockFacade: &MockMetricFacade{},
			expectedMetrics: []*domain.Metric{
				{
					ID:    "metric1",
					Value: "100",
				},
			},
			collectCalled:   true,
			sendCalled:      true,
			interruptSignal: false,
		},
		{
			name: "should stop on interrupt signal",
			mockCollectors: []*MockMetricsCollector{
				&MockMetricsCollector{},
			},
			mockFacade:      &MockMetricFacade{},
			expectedMetrics: nil,
			collectCalled:   false,
			sendCalled:      false,
			interruptSignal: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокаем Collect
			for _, collector := range tt.mockCollectors {
				collector.On("Collect").Return(tt.expectedMetrics)
			}

			// Мокаем UpdateMetric
			tt.mockFacade.On("UpdateMetric", mock.Anything).Return(nil)

			// Преобразуем коллекционеров в интерфейс MetricsCollector
			var collectors []MetricsCollector
			for _, collector := range tt.mockCollectors {
				collectors = append(collectors, collector)
			}

			// Создаем конфигурацию и агент
			config := &configs.AgentConfig{
				Address:        "localhost:8080",
				PollInterval:   100 * time.Millisecond,
				ReportInterval: 200 * time.Millisecond,
			}
			logger, _ := zap.NewProduction()
			agent := NewMetricAgent(config, collectors, tt.mockFacade, logger)

			// Канал для сигнала завершения
			signalCh := make(chan os.Signal, 1)
			if tt.interruptSignal {
				go func() {
					time.Sleep(50 * time.Millisecond)
					signalCh <- os.Interrupt
				}()
			}

			// Запускаем агент в отдельной горутине
			go agent.Run(signalCh)

			// Ждем некоторое время для имитации работы
			time.Sleep(300 * time.Millisecond)

			// Проверяем, были ли вызваны методы
			if tt.collectCalled {
				for _, collector := range tt.mockCollectors {
					collector.AssertExpectations(t)
				}
			}
			if tt.sendCalled {
				tt.mockFacade.AssertExpectations(t)
			}
		})
	}
}
