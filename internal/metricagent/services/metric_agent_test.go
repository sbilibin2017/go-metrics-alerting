package services

import (
	"errors"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для MetricFacade
type MockMetricFacade struct {
	mock.Mock
}

func (m *MockMetricFacade) UpdateMetric(metricType string, metricName string, metricValue string) error {
	args := m.Called(metricType, metricName, metricValue)
	return args.Error(0)
}

// Тестирование функционала MetricAgentService
func TestMetricAgentService_Start(t *testing.T) {
	// Создаем моки для facade
	mockMetricFacade := new(MockMetricFacade)

	// Настройка моки для обоих типов метрик
	mockMetricFacade.On("UpdateMetric", "gauge", mock.Anything, mock.Anything).Return(nil)
	mockMetricFacade.On("UpdateMetric", "counter", mock.Anything, mock.Anything).Return(nil)

	// Создаем сервис с моками
	service := NewMetricAgentService(mockMetricFacade, 1*time.Second, 1*time.Second)

	// Запускаем сервис в горутине, чтобы он не блокировал тест
	go service.Start()

	// Подождем немного, чтобы сервис собрал и отправил метрики
	time.Sleep(2 * time.Second)

	// Проверим, что метод был вызван хотя бы один раз с ожидаемыми аргументами
	mockMetricFacade.AssertExpectations(t)
}

func TestMetricAgentService_StartWithError(t *testing.T) {
	// Создаем мок
	mockFacade := new(MockMetricFacade)

	// Создаем экземпляр MetricAgentService с моками
	service := NewMetricAgentService(mockFacade, 10*time.Millisecond, 10*time.Millisecond)

	// Ожидаем, что UpdateMetric вызовет ошибку для "gauge" и метрики "Alloc"
	mockFacade.On("UpdateMetric", "gauge", "Alloc", mock.Anything).Return(errors.New("update error")).Times(1)
	// Замокируем все остальные метрики с помощью mock.Anything для значений
	mockFacade.On("UpdateMetric", "gauge", "BuckHashSys", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "Frees", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "HeapAlloc", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "HeapIdle", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "HeapInuse", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "HeapObjects", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "HeapReleased", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "HeapSys", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "NumGC", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "Sys", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "TotalAlloc", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "RandomValue", mock.Anything).Return(nil).Times(1)
	mockFacade.On("UpdateMetric", "gauge", "GCCPUFraction", mock.Anything).Return(nil).Times(1)

	// Для метрики типа "counter"
	mockFacade.On("UpdateMetric", "counter", "PollCount", mock.Anything).Return(nil).Times(1)

	// Запускаем сервис в отдельной горутине
	go func() {
		service.Start()
	}()

	// Даем сервису немного времени для работы
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что вызовы методов были сделаны
	mockFacade.AssertExpectations(t)
}

func TestMetricAgentService_Shutdown(t *testing.T) {
	// Создаем мок
	mockFacade := new(MockMetricFacade)

	// Создаем экземпляр MetricAgentService с моками
	service := NewMetricAgentService(mockFacade, 10*time.Millisecond, 10*time.Millisecond)

	// Запускаем сервис в отдельной горутине
	go func() {
		service.Start()
	}()

	// Отправляем сигнал SIGTERM для завершения работы
	service.shutdown <- syscall.SIGTERM

	// Даем сервису время завершить работу
	time.Sleep(50 * time.Millisecond)

	// Проверяем, что сервис завершил работу (например, по логам)
	// В тестах лучше всего использовать mock-объекты для проверки взаимодействий
	mockFacade.AssertExpectations(t)
}

func TestCollectGaugeMetrics(t *testing.T) {
	// Проверим, что collectGaugeMetrics собирает хотя бы одну метрику
	metrics := collectGaugeMetrics()

	// Проверяем, что собралось несколько метрик
	assert.Greater(t, len(metrics), 0)

	// Пример проверки, что хотя бы одна метрика правильная
	assert.Equal(t, "gauge", string(metrics[0].Type)) // Приводим MetricType к строке
	assert.Equal(t, "Alloc", metrics[0].Name)
	assert.NotEmpty(t, metrics[0].Value)
}

func TestCollectCounterMetrics(t *testing.T) {
	// Проверим, что collectCounterMetrics собирает хотя бы одну метрику
	metrics := collectCounterMetrics()

	// Проверяем, что собралось хотя бы одна метрика
	assert.Greater(t, len(metrics), 0)

	// Проверяем, что метрика с правильным именем и значением
	assert.Equal(t, "PollCount", metrics[0].Name)
	assert.Equal(t, "counter", string(metrics[0].Type)) // Сравниваем через тип MetricType
	assert.NotEmpty(t, metrics[0].Value)
}
