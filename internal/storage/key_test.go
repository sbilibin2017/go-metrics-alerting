package storage

import (
	"testing"
)

func TestEncode(t *testing.T) {
	km := NewKeyProcessor()

	tests := []struct {
		metricType string
		metricName string
		expected   string
	}{
		{"gauge", "cpu", "gauge:cpu"},
		{"counter", "requests", "counter:requests"},
		{"timer", "response_time", "timer:response_time"},
	}

	for _, test := range tests {
		t.Run(test.metricType+"_"+test.metricName, func(t *testing.T) {
			got := km.Encode(test.metricType, test.metricName)
			if got != test.expected {
				t.Errorf("Encode(%v, %v) = %v; want %v", test.metricType, test.metricName, got, test.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	km := NewKeyProcessor()

	tests := []struct {
		key          string
		expectedType string
		expectedName string
		expectedErr  error
	}{
		{"gauge:cpu", "gauge", "cpu", nil},
		{"counter:requests", "counter", "requests", nil},
		{"timer:response_time", "timer", "response_time", nil},
		{"invalidKey", "", "", ErrInvalidKeyFormat},
		{"gauge:", "", "", ErrInvalidKeyFormat},
		{":cpu", "", "", ErrInvalidKeyFormat},
	}

	for _, test := range tests {
		t.Run(test.key, func(t *testing.T) {
			gotType, gotName, err := km.Decode(test.key)

			if err != nil && err != test.expectedErr {
				t.Errorf("Decode(%v) returned error %v, want %v", test.key, err, test.expectedErr)
			}
			if gotType != test.expectedType {
				t.Errorf("Decode(%v) = %v, want %v", test.key, gotType, test.expectedType)
			}
			if gotName != test.expectedName {
				t.Errorf("Decode(%v) = %v, want %v", test.key, gotName, test.expectedName)
			}
		})
	}
}
