package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для интерфейса GetAllRepo
type MockGetAllRepo struct {
	mock.Mock
}

func (m *MockGetAllRepo) GetAll() [][3]string {
	args := m.Called()
	return args.Get(0).([][3]string)
}

func TestGetAllMetricValuesService_GetAllMetricValues_Success(t *testing.T) {
	mockRepo := new(MockGetAllRepo)
	service := NewGetAllMetricsService(mockRepo)

	// Мокируем возвращаемые значения для метрик
	mockRepo.On("GetAll").Return([][3]string{
		{"gauge", "cpu", "99"},
		{"counter", "requests", "1000"},
	})

	// Тестируем получение всех метрик
	metrics := service.GetAllMetricValues()

	assert.Len(t, metrics, 2) // Проверяем, что вернулось две метрики

	// Проверяем содержание метрик
	assert.Equal(t, "gauge", metrics[0].Type)
	assert.Equal(t, "cpu", metrics[0].Name)
	assert.Equal(t, "99", metrics[0].Value)

	assert.Equal(t, "counter", metrics[1].Type)
	assert.Equal(t, "requests", metrics[1].Name)
	assert.Equal(t, "1000", metrics[1].Value)

	mockRepo.AssertExpectations(t)
}

func TestGetAllMetricValuesService_GetAllMetricValues_Empty(t *testing.T) {
	mockRepo := new(MockGetAllRepo)
	service := NewGetAllMetricsService(mockRepo)

	// Мокируем возвращаемое пустое значение (отсутствие метрик)
	mockRepo.On("GetAll").Return([][3]string{})

	// Тестируем случай, когда метрики отсутствуют
	metrics := service.GetAllMetricValues()

	assert.Len(t, metrics, 0) // Проверяем, что вернулся пустой срез
	mockRepo.AssertExpectations(t)
}
