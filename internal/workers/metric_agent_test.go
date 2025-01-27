package workers_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sbilibin2017/go-metrics-alerting/internal/workers"
	"github.com/stretchr/testify/mock"
)

// Мок для MetricAgentServiceInterface
// Мок для MetricAgentServiceInterface
type MockAgentService struct {
	mock.Mock
}

func (m *MockAgentService) CollectMetrics() error {
	args := m.Called()
	return args.Error(0) // Возвращаем ошибку, если она была передана в вызов
}

// Мок для AgentConfigInterface
type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) GetPollInterval() int {
	args := m.Called()
	return args.Int(0)
}

func (m *MockConfig) GetReportInterval() int {
	args := m.Called()
	return args.Int(0)
}

// Добавим реализацию метода GetBaseUrl
func (m *MockConfig) GetBaseUrl() string {
	args := m.Called()
	return args.String(0)
}

func TestMetricsWorker_Start_SuccessfulMethodCalls(t *testing.T) {
	mockAgentService := new(MockAgentService)
	mockConfig := new(MockConfig)

	// Настроим моки
	mockConfig.On("GetPollInterval").Return(1)
	mockConfig.On("GetReportInterval").Return(1)
	mockConfig.On("GetBaseUrl").Return("http://example.com")

	// Настроим мок для CollectMetrics
	mockAgentService.On("CollectMetrics").Return(nil)

	// Создаем рабочий объект с моками
	worker := workers.NewMetricsWorker(mockAgentService, mockConfig)

	// Тестируем Start
	go worker.Start()

	// Даем немного времени, чтобы потикеры могли сработать
	time.Sleep(2 * time.Second)

	// Проверяем, что метод CollectMetrics был вызван хотя бы один раз.
	mockAgentService.AssertCalled(t, "CollectMetrics")

	// Проверяем, что методы GetPollInterval, GetReportInterval и GetBaseUrl были вызваны.
	mockConfig.AssertCalled(t, "GetPollInterval")
	mockConfig.AssertCalled(t, "GetReportInterval")
	mockConfig.AssertCalled(t, "GetBaseUrl")
}

func TestMetricsWorker_Start_CollectMetricsHandlesError(t *testing.T) {
	mockAgentService := new(MockAgentService)
	mockConfig := new(MockConfig)

	// Настроим моки
	mockConfig.On("GetPollInterval").Return(1)
	mockConfig.On("GetReportInterval").Return(1)
	mockConfig.On("GetBaseUrl").Return("http://example.com") // добавляем настройку для GetBaseUrl

	// Настроим мок для CollectMetrics, чтобы он возвращал ошибку
	mockAgentService.On("CollectMetrics").Return(fmt.Errorf("some error"))

	// Создаем рабочий объект с моками
	worker := workers.NewMetricsWorker(mockAgentService, mockConfig)

	// Тестируем Start
	go worker.Start()

	// Даем немного времени, чтобы потикеры могли сработать
	time.Sleep(2 * time.Second)

	// Проверяем, что метод CollectMetrics был вызван хотя бы один раз, несмотря на ошибку
	mockAgentService.AssertCalled(t, "CollectMetrics")
}

func TestMetricsWorker_Start_ZeroInterval(t *testing.T) {
	mockAgentService := new(MockAgentService)
	mockConfig := new(MockConfig)

	// Настроим моки с минимальными интервалами (например, 1 секунда)
	mockConfig.On("GetPollInterval").Return(1)
	mockConfig.On("GetReportInterval").Return(1)
	mockConfig.On("GetBaseUrl").Return("http://example.com") // добавляем настройку для GetBaseUrl

	// Настроим мок для CollectMetrics
	mockAgentService.On("CollectMetrics").Return(nil)

	// Создаем рабочий объект с моками
	worker := workers.NewMetricsWorker(mockAgentService, mockConfig)

	// Тестируем Start
	go worker.Start()

	// Даем немного времени, чтобы потикеры могли сработать
	time.Sleep(2 * time.Second)

	// Проверяем, что метод CollectMetrics был вызван хотя бы один раз.
	mockAgentService.AssertCalled(t, "CollectMetrics")

	// Проверяем, что методы GetPollInterval и GetReportInterval были вызваны.
	mockConfig.AssertCalled(t, "GetPollInterval")
	mockConfig.AssertCalled(t, "GetReportInterval")
}

func TestMetricsWorker_Start_LongIntervals(t *testing.T) {
	mockAgentService := new(MockAgentService)
	mockConfig := new(MockConfig)

	// Настроим моки с интервалами
	mockConfig.On("GetPollInterval").Return(10)              // Например, 10 секунд
	mockConfig.On("GetReportInterval").Return(20)            // Например, 20 секунд
	mockConfig.On("GetBaseUrl").Return("http://example.com") // добавляем настройку для GetBaseUrl

	// Настроим мок для CollectMetrics
	mockAgentService.On("CollectMetrics").Return(nil)

	// Создаем рабочий объект с моками
	worker := workers.NewMetricsWorker(mockAgentService, mockConfig)

	// Тестируем Start
	go worker.Start()

	// Даем немного времени, чтобы потикеры могли сработать
	time.Sleep(2 * time.Second)

	// Проверяем, что метод CollectMetrics был вызван хотя бы один раз.
	mockAgentService.AssertCalled(t, "CollectMetrics")

	// Проверяем, что методы GetPollInterval и GetReportInterval были вызваны.
	mockConfig.AssertCalled(t, "GetPollInterval")
	mockConfig.AssertCalled(t, "GetReportInterval")
	mockConfig.AssertCalled(t, "GetBaseUrl")
}
