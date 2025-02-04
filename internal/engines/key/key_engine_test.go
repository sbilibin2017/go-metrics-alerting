package key

import (
	"errors"
	"testing"
)

// Тест Encode
func TestKeyEngine_Encode(t *testing.T) {
	engine := NewKeyEngine()

	tests := []struct {
		name    string
		key     *Key
		want    string
		wantErr error
	}{
		{"Valid Key", &Key{"cpu", "usage"}, "cpu:usage", nil},
		{"Empty Type", &Key{"", "usage"}, "", ErrInvalidKeyFormat},
		{"Empty Name", &Key{"cpu", ""}, "", ErrInvalidKeyFormat},
		{"Nil Key", nil, "", ErrInvalidKeyFormat},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Encode(tt.key)

			if got != tt.want {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Тест Decode
func TestKeyEngine_Decode(t *testing.T) {
	engine := NewKeyEngine()

	tests := []struct {
		name    string
		key     string
		want    *Key
		wantErr error
	}{
		{"Valid Key", "cpu:usage", &Key{"cpu", "usage"}, nil},
		{"No Separator", "cpuusage", nil, ErrInvalidKeyFormat},
		{"Empty Type", ":usage", nil, ErrInvalidKeyFormat},
		{"Empty Name", "cpu:", nil, ErrInvalidKeyFormat},
		{"Double Separator", "cpu:usage:extra", &Key{"cpu", "usage:extra"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Decode(tt.key)

			if tt.want != nil {
				if got == nil || got.MetricType != tt.want.MetricType || got.MetricName != tt.want.MetricName {
					t.Errorf("Decode() = %+v, want %+v", got, tt.want)
				}
			} else if got != nil {
				t.Errorf("Decode() expected nil, got %+v", got)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
