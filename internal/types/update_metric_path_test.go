package types

import (
	"go-metrics-alerting/internal/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateMetricPathRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricPathRequest
		expected int
	}{
		{
			name: "Valid request with gauge and valid value",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: string(domain.Gauge),
				Value: "123.45", // Пример значения для gauge
			},
			expected: http.StatusOK, // Ожидаем, что ошибки не будет
		},
		{
			name: "Valid request with counter and valid value",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: string(domain.Counter),
				Value: "123", // Пример значения для counter
			},
			expected: http.StatusOK, // Ожидаем, что ошибки не будет
		},
		{
			name: "Invalid request with empty ID",
			request: UpdateMetricPathRequest{
				ID:    "",
				MType: string(domain.Gauge),
				Value: "123.45",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid request with invalid ID format",
			request: UpdateMetricPathRequest{
				ID:    "invalid id!",
				MType: string(domain.Gauge),
				Value: "123.45",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid request with invalid metric type",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: "invalid_type", // Некорректный тип метрики
				Value: "123.45",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid request with invalid value for gauge",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: string(domain.Gauge),
				Value: "not_a_number", // Некорректное значение для gauge
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid request with invalid value for counter",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: string(domain.Counter),
				Value: "not_a_number", // Некорректное значение для counter
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			// Проверка только статуса ошибки
			if tt.expected == http.StatusOK {
				// Если ожидается nil, то ошибка не должна быть
				assert.Nil(t, err)
			} else {
				// Если ожидается ошибка, то проверяем только статус
				assert.NotNil(t, err)
				assert.Equal(t, tt.expected, err.Status)
			}
		})
	}
}

func TestPathToDomain(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricPathRequest
		expected *domain.Metrics
	}{
		{
			name: "Valid request with gauge",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: "gauge",
				Value: "123.45",
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Gauge,
				Delta: nil,
				Value: floatPointer(123.45), // Вспомогательная функция для указателя на float64
			},
		},
		{
			name: "Valid request with counter",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: "counter",
				Value: "42",
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Counter,
				Delta: intPointer(42), // Вспомогательная функция для указателя на int64
				Value: nil,
			},
		},
		{
			name: "Invalid request with non-numeric value for gauge",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: "gauge",
				Value: "not_a_number",
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Gauge,
				Delta: nil,
				Value: nil, // Значение не должно быть установлено
			},
		},
		{
			name: "Invalid request with non-numeric value for counter",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: "counter",
				Value: "not_a_number",
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Counter,
				Delta: nil, // Дельта не должна быть установлена
				Value: nil,
			},
		},
		{
			name: "Request with missing value for gauge",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: "gauge",
				Value: "",
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Gauge,
				Delta: nil,
				Value: nil, // Если значение пустое, то должно быть nil
			},
		},
		{
			name: "Request with missing value for counter",
			request: UpdateMetricPathRequest{
				ID:    "valid_id",
				MType: "counter",
				Value: "",
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Counter,
				Delta: nil, // Если значение пустое, то должно быть nil
				Value: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act: Convert the request to domain metrics
			result := tt.request.ToDomain()

			// Assert: Compare the result with expected domain metrics
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Вспомогательная функция для получения указателя на float64
func floatPointer(f float64) *float64 {
	return &f
}

// Вспомогательная функция для получения указателя на int64
func intPointer(i int64) *int64 {
	return &i
}
