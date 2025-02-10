package services

import (
	"context"
	"go-metrics-alerting/internal/types"

	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Мок для интерфейса GetAllRepo
type MockGetAllRepo struct {
	mock.Mock
}

func (m *MockGetAllRepo) GetAll(ctx context.Context) [][]string {
	args := m.Called(ctx)
	return args.Get(0).([][]string)
}

func TestGetAllMetricValues(t *testing.T) {
	tests := []struct {
		name          string
		mockData      [][]string
		expected      []*types.MetricResponse
		expectedError bool
	}{
		{
			name:          "с пустым списком метрик",
			mockData:      [][]string{},
			expected:      []*types.MetricResponse{}, // ожидаем пустой срез
			expectedError: false,
		},
		{
			name: "с одним элементом в списке",
			mockData: [][]string{
				{"1", "metric1", "100"},
			},
			expected: []*types.MetricResponse{
				{
					Name:  "metric1",
					Value: "100",
				},
			},
			expectedError: false,
		},
		{
			name: "с несколькими элементами в списке",
			mockData: [][]string{
				{"1", "metric1", "100"},
				{"2", "metric2", "200"},
			},
			expected: []*types.MetricResponse{
				{
					Name:  "metric1",
					Value: "100",
				},
				{
					Name:  "metric2",
					Value: "200",
				},
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем мок-объект для репозитория
			mockRepo := new(MockGetAllRepo)
			// Устанавливаем поведение мока
			mockRepo.On("GetAll", mock.Anything).Return(tt.mockData)

			// Создаем сервис
			service := &GetAllMetricValuesService{
				MetricRepository: mockRepo,
			}

			// Вызываем тестируемую функцию
			result := service.GetAllMetricValues(context.Background())

			// Проверяем результаты
			assert.Equal(t, tt.expected, result)

			// Проверяем, что метод GetAll был вызван
			mockRepo.AssertExpectations(t)
		})
	}
}
