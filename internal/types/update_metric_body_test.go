package types

import (
	"go-metrics-alerting/internal/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUpdateMetricBodyRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricBodyRequest
		expected int
	}{
		{
			name: "Valid request with gauge",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "gauge",
				Value: new(float64), // Пример правильного значения для gauge
			},
			expected: http.StatusOK, // ожидаем, что ошибки не будет
		},
		{
			name: "Valid request with counter",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "counter",
				Delta: new(int64), // Пример правильного значения для counter
			},
			expected: http.StatusOK, // ожидаем, что ошибки не будет
		},
		{
			name: "Invalid request with empty ID",
			request: UpdateMetricBodyRequest{
				ID:    "",
				MType: "gauge",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid request with invalid ID format",
			request: UpdateMetricBodyRequest{
				ID:    "invalid id!",
				MType: "gauge",
			},
			expected: http.StatusNotFound,
		},
		{
			name: "Invalid request with invalid metric type",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "invalid_type",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid request with missing value for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "gauge",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid request with missing delta for counter",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "counter",
			},
			expected: http.StatusBadRequest,
		},
		{
			name: "Invalid request with both value and delta for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "gauge",
				Value: new(float64), // Значение для gauge
				Delta: new(int64),   // И delta, что недопустимо
			},
			expected: http.StatusBadRequest, // Ожидаем ошибку с кодом BadRequest
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

func TestBodyToDomain(t *testing.T) {
	tests := []struct {
		name     string
		request  UpdateMetricBodyRequest
		expected *domain.Metrics
	}{
		{
			name: "Valid request with gauge",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "gauge",
				Value: new(float64),
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Gauge,
				Delta: nil,
				Value: new(float64),
			},
		},
		{
			name: "Valid request with counter",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "counter",
				Delta: new(int64),
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Counter,
				Delta: new(int64),
				Value: nil,
			},
		},
		{
			name: "Request with no delta for counter",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "counter",
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Counter,
				Delta: nil,
				Value: nil,
			},
		},
		{
			name: "Request with no value for gauge",
			request: UpdateMetricBodyRequest{
				ID:    "valid_id",
				MType: "gauge",
			},
			expected: &domain.Metrics{
				ID:    "valid_id",
				MType: domain.Gauge,
				Delta: nil,
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
