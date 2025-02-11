package engines

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyEngine_Encode(t *testing.T) {
	ke := &KeyEngine{}

	tests := []struct {
		name     string
		mt       string
		mn       string
		expected string
	}{
		{"valid input", "gauge", "cpu", "gauge:cpu"},
		{"empty metric name", "counter", "", "counter:"},
		{"empty metric type", "", "memory", ":memory"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ke.Encode(tt.mt, tt.mn)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestKeyEngine_Decode(t *testing.T) {
	ke := &KeyEngine{}

	tests := []struct {
		name   string
		key    string
		expMt  string
		expMn  string
		expErr error
	}{
		{"valid key", "gauge:cpu", "gauge", "cpu", nil},
		{"invalid key format - missing separator", "gaugecpu", "", "", ErrInvalidKeyFormat},
		{"invalid key format - empty type", ":memory", "", "", ErrInvalidKeyFormat},
		{"invalid key format - empty name", "counter:", "", "", ErrInvalidKeyFormat},
		{"invalid key format - multiple separators", "gauge:cpu:extra", "", "", ErrInvalidKeyFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mt, mn, err := ke.Decode(tt.key)
			assert.Equal(t, tt.expMt, mt)
			assert.Equal(t, tt.expMn, mn)
			assert.Equal(t, tt.expErr, err)
		})
	}
}
