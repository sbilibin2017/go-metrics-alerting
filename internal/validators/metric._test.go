package validators_test

import (
	"testing"

	"go-metrics-alerting/internal/types"
	"go-metrics-alerting/internal/validators"

	"github.com/stretchr/testify/assert"
)

func TestValidateEmptyString(t *testing.T) {
	v := &validators.ValidateEmptyString{}

	// Test when ID is empty
	err := v.Validate("")
	assert.Equal(t, err, validators.ErrEmptyID)

	// Test when ID is not empty
	err = v.Validate("some_id")
	assert.NoError(t, err)
}

func TestValidateMType(t *testing.T) {
	v := &validators.ValidateMType{}

	// Test invalid metric type
	err := v.Validate("invalid_type")
	assert.Equal(t, err, validators.ErrInvalidMType)

	// Test valid metric types
	err = v.Validate(types.Counter)
	assert.NoError(t, err)

	err = v.Validate(types.Gauge)
	assert.NoError(t, err)
}

func TestValidateDelta(t *testing.T) {
	v := &validators.ValidateDelta{}

	// Test when mType is Counter and delta is nil
	err := v.Validate(types.Counter, nil)
	assert.Equal(t, err, validators.ErrInvalidDelta)

	// Test when mType is Counter and delta is not nil
	delta := int64(10)
	err = v.Validate(types.Counter, &delta)
	assert.NoError(t, err)

	// Test when mType is Gauge and delta is nil (should not return error)
	err = v.Validate(types.Gauge, nil)
	assert.NoError(t, err)
}

func TestValidateValue(t *testing.T) {
	v := &validators.ValidateValue{}

	// Test when mType is Gauge and value is nil
	err := v.Validate(types.Gauge, nil)
	assert.Equal(t, err, validators.ErrInvalidValue)

	// Test when mType is Gauge and value is not nil
	value := float64(10.5)
	err = v.Validate(types.Gauge, &value)
	assert.NoError(t, err)

	// Test when mType is Counter and value is nil (should not return error)
	err = v.Validate(types.Counter, nil)
	assert.NoError(t, err)
}

func TestValidateCounterValue(t *testing.T) {
	v := &validators.ValidateCounterValue{}

	// Test invalid counter value
	err := v.Validate("invalid_value")
	assert.Equal(t, err, validators.ErrInvalidCounterVal)

	// Test valid counter value
	err = v.Validate("100")
	assert.NoError(t, err)
}

func TestValidateGaugeValue(t *testing.T) {
	v := &validators.ValidateGaugeValue{}

	// Test invalid gauge value
	err := v.Validate("invalid_value")
	assert.Equal(t, err, validators.ErrInvalidGaugeVal)

	// Test valid gauge value
	err = v.Validate("100.5")
	assert.NoError(t, err)
}
