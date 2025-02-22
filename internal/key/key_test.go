package key

import (
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		values   []string
		expected string
	}{
		{
			name:     "Single value",
			values:   []string{"metric"},
			expected: "metric",
		},
		{
			name:     "Two values",
			values:   []string{"metric", "12345"},
			expected: "metric:12345",
		},
		{
			name:     "Multiple values",
			values:   []string{"metric", "12345", "extra_value"},
			expected: "metric:12345:extra_value",
		},
		{
			name:     "Empty values",
			values:   []string{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Encode(tt.values...)
			if got != tt.expected {
				t.Errorf("Encode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected []string
	}{
		{
			name:     "Single value",
			key:      "metric",
			expected: []string{"metric"},
		},
		{
			name:     "Two values",
			key:      "metric:12345",
			expected: []string{"metric", "12345"},
		},
		{
			name:     "Multiple values",
			key:      "metric:12345:extra_value",
			expected: []string{"metric", "12345", "extra_value"},
		},
		{
			name:     "Empty string",
			key:      "",
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Decode(tt.key)
			if !equal(got, tt.expected) {
				t.Errorf("Decode() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// Вспомогательная функция для сравнения срезов строк
func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
