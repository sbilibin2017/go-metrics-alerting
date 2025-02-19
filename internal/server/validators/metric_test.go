package validators

import (
	"go-metrics-alerting/internal/server/types"
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
			err := ValidateEmptyString(tt.id)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateMType(t *testing.T) {
	tests := []struct {
		mType       types.MType
		expectedErr error
	}{
		{types.Counter, nil},              // Тип Counter должен быть валидным
		{types.Gauge, nil},                // Тип Gauge должен быть валидным
		{"invalid-type", ErrInvalidMType}, // Неверный тип метрики должен привести к ошибке
	}

	for _, tt := range tests {
		t.Run(string(tt.mType), func(t *testing.T) {
			err := ValidateMType(tt.mType)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateDelta(t *testing.T) {
	tests := []struct {
		mType       types.MType
		delta       *int64
		expectedErr error
	}{
		{types.Counter, nil, ErrInvalidDelta}, // Для Counter, если Delta = nil, ошибка
		{types.Counter, new(int64), nil},      // Для Counter, если Delta задано, ошибок нет
		{types.Gauge, nil, nil},               // Для Gauge нет необходимости в Delta, ошибок нет
	}

	for _, tt := range tests {
		t.Run(string(tt.mType), func(t *testing.T) {
			err := ValidateDelta(tt.mType, tt.delta)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateValue(t *testing.T) {
	tests := []struct {
		mType       types.MType
		value       *float64
		expectedErr error
	}{
		{types.Gauge, nil, ErrInvalidValue}, // Для Gauge, если Value = nil, ошибка
		{types.Gauge, new(float64), nil},    // Для Gauge, если Value задано, ошибок нет
		{types.Counter, nil, nil},           // Для Counter нет необходимости в Value, ошибок нет
	}

	for _, tt := range tests {
		t.Run(string(tt.mType), func(t *testing.T) {
			err := ValidateValue(tt.mType, tt.value)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
