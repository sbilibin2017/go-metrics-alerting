package engines

import (
	"testing"
)

func TestKeyEngine_Encode(t *testing.T) {
	tests := []struct {
		name       string
		metricType string
		metricName string
		expected   string
	}{
		{
			name:       "Valid encoding",
			metricType: "counter",
			metricName: "requests",
			expected:   "counter:requests",
		},
		{
			name:       "Another valid encoding",
			metricType: "gauge",
			metricName: "temperature",
			expected:   "gauge:temperature",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeyEngine()
			result := k.Encode(tt.metricType, tt.metricName)

			if result != tt.expected {
				t.Errorf("KeyEngine.Encode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestKeyEngine_Decode(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		expectedType string
		expectedName string
		expectError  bool
	}{
		{
			name:         "Valid decoding",
			key:          "counter:requests",
			expectedType: "counter",
			expectedName: "requests",
			expectError:  false,
		},
		{
			name:         "Another valid decoding",
			key:          "gauge:temperature",
			expectedType: "gauge",
			expectedName: "temperature",
			expectError:  false,
		},
		{
			name:         "Invalid key format (too few parts)",
			key:          "counter",
			expectedType: "",
			expectedName: "",
			expectError:  true,
		},
		{
			name:         "Invalid key format (too many parts)",
			key:          "counter:requests:extra",
			expectedType: "",
			expectedName: "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := NewKeyEngine()
			metricType, metricName, err := k.Decode(tt.key)

			if (err != nil) != tt.expectError {
				t.Errorf("KeyEngine.Decode() error = %v, wantErr %v", err, tt.expectError)
				return
			}

			if err == nil {
				if metricType != tt.expectedType {
					t.Errorf("KeyEngine.Decode() metricType = %v, want %v", metricType, tt.expectedType)
				}
				if metricName != tt.expectedName {
					t.Errorf("KeyEngine.Decode() metricName = %v, want %v", metricName, tt.expectedName)
				}
			}
		})
	}
}
