package engines

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyEngine_Encode(t *testing.T) {
	tests := []struct {
		metricType string
		metricName string
		expected   string
	}{
		{"counter", "metric1", "counter:metric1"},
		{"gauge", "metric2", "gauge:metric2"},
		{"counter", "metric3", "counter:metric3"},
	}

	ke := &KeyEngine{}
	for _, test := range tests {
		t.Run(test.metricType+"_"+test.metricName, func(t *testing.T) {
			result := ke.Encode(test.metricType, test.metricName)
			assert.Equal(t, test.expected, result, "Encoded key should match expected value")
		})
	}
}

func TestKeyEngine_Decode(t *testing.T) {
	tests := []struct {
		key          string
		expectedType string
		expectedName string
		expectError  bool
	}{
		{"counter:metric1", "counter", "metric1", false},
		{"gauge:metric2", "gauge", "metric2", false},
		{"invalidkey", "", "", true},
		{":missingname", "", "", true},
		{"missingtype:", "", "", true},
	}

	ke := &KeyEngine{}
	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			metricType, metricName, err := ke.Decode(test.key)

			if test.expectError {
				require.Error(t, err, "Expected an error but got none")
			} else {
				require.NoError(t, err, "Unexpected error occurred")
				assert.Equal(t, test.expectedType, metricType, "Decoded metric type should match expected")
				assert.Equal(t, test.expectedName, metricName, "Decoded metric name should match expected")
			}
		})
	}
}

func TestKeyEngine_Decode_EmptyStrings(t *testing.T) {
	ke := &KeyEngine{}
	tests := []struct {
		key         string
		expectError bool
	}{
		{"", true},
		{":", true},
		{"emptykey:", true},
		{":emptykey", true},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			_, _, err := ke.Decode(test.key)
			if test.expectError {
				require.Error(t, err, "Expected an error for empty key but got none")
			} else {
				require.NoError(t, err, "Unexpected error for key")
			}
		})
	}
}
