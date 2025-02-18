package validators

import (
	"go-metrics-alerting/internal/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateID(t *testing.T) {
	tests := []struct {
		id          string
		expectedErr error
	}{
		{"", ErrEmptyID},  // Пустой ID должен привести к ошибке
		{"valid-id", nil}, // Валидный ID должен пройти без ошибок
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			err := ValidateID(tt.id)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateMType(t *testing.T) {
	tests := []struct {
		mType       string
		expectedErr error
	}{
		{string(types.Counter), nil},      // Тип Counter должен быть валидным
		{string(types.Gauge), nil},        // Тип Gauge должен быть валидным
		{"invalid-type", ErrInvalidMType}, // Неверный тип метрики должен привести к ошибке
	}

	for _, tt := range tests {
		t.Run(tt.mType, func(t *testing.T) {
			err := ValidateMType(tt.mType)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateDelta(t *testing.T) {
	tests := []struct {
		mType       string
		delta       *int64
		expectedErr error
	}{
		{string(types.Counter), nil, ErrInvalidDelta}, // Для Counter, если Delta = nil, ошибка
		{string(types.Counter), new(int64), nil},      // Для Counter, если Delta задано, ошибок нет
		{string(types.Gauge), nil, nil},               // Для Gauge нет необходимости в Delta, ошибок нет
	}

	for _, tt := range tests {
		t.Run(tt.mType, func(t *testing.T) {
			err := ValidateDelta(tt.mType, tt.delta)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateValue(t *testing.T) {
	tests := []struct {
		mType       string
		value       *float64
		expectedErr error
	}{
		{string(types.Gauge), nil, ErrInvalidValue}, // Для Gauge, если Value = nil, ошибка
		{string(types.Gauge), new(float64), nil},    // Для Gauge, если Value задано, ошибок нет
		{string(types.Counter), nil, nil},           // Для Counter нет необходимости в Value, ошибок нет
	}

	for _, tt := range tests {
		t.Run(tt.mType, func(t *testing.T) {
			err := ValidateValue(tt.mType, tt.value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
