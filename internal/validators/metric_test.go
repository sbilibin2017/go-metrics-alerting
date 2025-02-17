package validators

import (
	"go-metrics-alerting/internal/types"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDValidator_Valid(t *testing.T) {
	validator := &IDValidator{}

	// Testing for the expected valid case (EmptyString)
	result := validator.Validate("")
	assert.True(t, result, "IDValidator should return true for EmptyString")
}

func TestIDValidator_Invalid(t *testing.T) {
	validator := &IDValidator{}

	// Testing for an invalid ID (non-empty string)
	result := validator.Validate("non-empty-id")
	assert.False(t, result, "IDValidator should return false for non-empty ID")
}

func TestMTypeValidator_Valid(t *testing.T) {
	validator := &MTypeValidator{}

	// Testing for an invalid metric type (Counter and Gauge are valid, so should return false)
	resultCounter := validator.Validate(string(types.Counter))
	assert.False(t, resultCounter, "MTypeValidator should return false for Counter")

	resultGauge := validator.Validate(string(types.Gauge))
	assert.False(t, resultGauge, "MTypeValidator should return false for Gauge")
}

func TestMTypeValidator_Invalid(t *testing.T) {
	validator := &MTypeValidator{}

	// Testing for an invalid metric type (should return true for invalid types)
	resultInvalid := validator.Validate("invalid_type")
	assert.True(t, resultInvalid, "MTypeValidator should return true for an invalid type")
}

func TestDeltaValidator_Valid(t *testing.T) {
	validator := &DeltaValidator{}

	// Testing for the expected valid case (Counter metric with delta == nil)
	var delta *int64
	result := validator.Validate(string(types.Counter), delta)
	assert.True(t, result, "DeltaValidator should return true for Counter metric with nil delta")
}

func TestDeltaValidator_Invalid(t *testing.T) {
	validator := &DeltaValidator{}

	// Testing for invalid Counter metric with a non-nil delta
	delta := int64(5)
	result := validator.Validate(string(types.Counter), &delta)
	assert.False(t, result, "DeltaValidator should return false for non-nil delta for Counter metric")

	// Testing for non-Counter metric (should return false)
	resultNonCounter := validator.Validate("gauge", &delta)
	assert.False(t, resultNonCounter, "DeltaValidator should return false for non-Counter metric")
}

func TestValueValidator_Valid(t *testing.T) {
	validator := &ValueValidator{}

	// Testing for the expected valid case (Gauge metric with value == nil)
	var value *float64
	result := validator.Validate(string(types.Gauge), value)
	assert.True(t, result, "ValueValidator should return true for Gauge metric with nil value")
}

func TestValueValidator_Invalid(t *testing.T) {
	validator := &ValueValidator{}

	// Testing for invalid Gauge metric with a non-nil value
	value := float64(3.14)
	result := validator.Validate(string(types.Gauge), &value)
	assert.False(t, result, "ValueValidator should return false for non-nil value for Gauge metric")

	// Testing for non-Gauge metric (should return false)
	resultNonGauge := validator.Validate("counter", &value)
	assert.False(t, resultNonGauge, "ValueValidator should return false for non-Gauge metric")
}
