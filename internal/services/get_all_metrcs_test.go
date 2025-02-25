package services

import (
	"go-metrics-alerting/internal/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Создадим мок для интерфейса Ranger
type MockRanger struct {
	mock.Mock
}

func (m *MockRanger) Range(callback func(key string, value *domain.Metrics) bool) {
	args := m.Called(callback)
	if args.Bool(0) {
		// Используем указатели на float64 для Value
		value1 := 10.0
		value2 := 20.0
		callback("metric1", &domain.Metrics{ID: "1", Value: &value1})
		callback("metric2", &domain.Metrics{ID: "2", Value: &value2})
	}
}

// Тест для метода GetAllMetrics
func TestGetAllMetrics(t *testing.T) {
	// Создаем мок для Ranger
	mockRanger := new(MockRanger)

	// Создаем экземпляр GetAllMetricsService с мокированным Ranger
	service := NewGetAllMetricsService(mockRanger)

	// Настроим ожидаемое поведение для Range
	mockRanger.On("Range", mock.Anything).Return(true).Once()

	// Вызовем метод GetAllMetrics
	metrics := service.GetAllMetrics()

	// Проверим, что метрики были правильно получены
	assert.Equal(t, 2, len(metrics), "expected 2 metrics")
	assert.Equal(t, "1", metrics[0].ID, "expected first metric ID to be 1")
	assert.Equal(t, 10.0, *metrics[0].Value, "expected first metric value to be 10")
	assert.Equal(t, "2", metrics[1].ID, "expected second metric ID to be 2")
	assert.Equal(t, 20.0, *metrics[1].Value, "expected second metric value to be 20")

	// Убедимся, что метод Range был вызван
	mockRanger.AssertExpectations(t)
}
