package engines

import (
	"testing"
)

func TestCounterStrategy_Update(t *testing.T) {
	tests := []struct {
		name         string
		currentValue string
		newValue     string
		expected     string
		expectError  bool
	}{
		{
			name:         "Valid increment",
			currentValue: "10",
			newValue:     "5",
			expected:     "15",
			expectError:  false,
		},
		{
			name:         "Decrement to zero",
			currentValue: "10",
			newValue:     "-10",
			expected:     "0",
			expectError:  false,
		},
		{
			name:         "Decrement to negative",
			currentValue: "10",
			newValue:     "-15",
			expected:     "-5",
			expectError:  false,
		},
		{
			name:         "Negative current value",
			currentValue: "-10",
			newValue:     "5",
			expected:     "-5",
			expectError:  false,
		},
		{
			name:         "Invalid current value",
			currentValue: "abc",
			newValue:     "5",
			expected:     "",
			expectError:  true,
		},
		{
			name:         "Invalid increment value",
			currentValue: "10",
			newValue:     "xyz",
			expected:     "",
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &CounterUpdateStrategyEngine{}
			result, err := cs.Update(tt.currentValue, tt.newValue)

			if (err != nil) != tt.expectError {
				t.Errorf("CounterStrategy.Update() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if result != tt.expected {
				t.Errorf("CounterStrategy.Update() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGaugeStrategy_Update(t *testing.T) {
	tests := []struct {
		name         string
		currentValue string
		newValue     string
		expected     string
		expectError  bool
	}{
		{
			name:         "Valid gauge value",
			currentValue: "10.5",
			newValue:     "25.75",
			expected:     "25.75",
			expectError:  false,
		},
		{
			name:         "Invalid new value (non-numeric)",
			currentValue: "10.5",
			newValue:     "abc",
			expected:     "",
			expectError:  true,
		},
		{
			name:         "Invalid new value (empty string)",
			currentValue: "10.5",
			newValue:     "",
			expected:     "",
			expectError:  true,
		},
		{
			name:         "Negative gauge value",
			currentValue: "10.5",
			newValue:     "-5.5",
			expected:     "-5.5",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gs := &GaugeUpdateStrategyEngine{}
			result, err := gs.Update(tt.currentValue, tt.newValue)
			if (err != nil) != tt.expectError {
				t.Errorf("GaugeStrategy.Update() error = %v, wantErr %v", err, tt.expectError)
				return
			}
			if result != tt.expected {
				t.Errorf("GaugeStrategy.Update() = %v, want %v", result, tt.expected)
			}
		})
	}
}
