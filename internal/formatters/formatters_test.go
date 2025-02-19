package formatters_test

import (
	"testing"

	"go-metrics-alerting/internal/formatters"

	"github.com/stretchr/testify/assert"
)

func TestInt64Handler_Parse(t *testing.T) {
	h := &formatters.Int64Formatter{}

	// Test valid int64 value
	value, err := h.Parse("12345")
	assert.NoError(t, err)
	assert.Equal(t, int64(12345), value)

	// Test invalid int64 value
	_, err = h.Parse("invalid")
	assert.EqualError(t, err, "failed to parse int64: strconv.ParseInt: parsing \"invalid\": invalid syntax")
}

func TestInt64Handler_Format(t *testing.T) {
	h := &formatters.Int64Formatter{}

	// Test formatting int64 value
	value := int64(12345)
	formattedValue := h.Format(value)
	assert.Equal(t, "12345", formattedValue)
}

func TestFloat64Handler_Parse(t *testing.T) {
	h := &formatters.Float64Formatter{}

	// Test valid float64 value
	value, err := h.Parse("123.45")
	assert.NoError(t, err)
	assert.Equal(t, 123.45, value)

	// Test invalid float64 value
	_, err = h.Parse("invalid")
	assert.EqualError(t, err, "failed to parse float64: strconv.ParseFloat: parsing \"invalid\": invalid syntax")
}

func TestFloat64Handler_Format(t *testing.T) {
	h := &formatters.Float64Formatter{}

	// Test formatting float64 value
	value := 123.45
	formattedValue := h.Format(value)
	assert.Equal(t, "123.45", formattedValue)

	// Test formatting float64 value with scientific notation
	value = 1.23e5
	formattedValue = h.Format(value)
	assert.Equal(t, "123000", formattedValue)
}
