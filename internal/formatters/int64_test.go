package formatters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Тест для FormatInt64
func TestFormatInt64(t *testing.T) {
	tests := []struct {
		value    int64
		expected string
	}{
		{123456, "123456"},
		{-123456, "-123456"},
		{0, "0"},
		{1000000000, "1000000000"},
	}

	for _, test := range tests {
		t.Run("Test FormatInt64", func(t *testing.T) {
			result := FormatInt64(test.value)
			assert.Equal(t, test.expected, result, "they should be equal")
		})
	}
}

// Тест для ParseInt64
func TestParseInt64(t *testing.T) {
	tests := []struct {
		value    string
		expected int64
		ok       bool // Обновлено: теперь проверяем значение `ok` вместо `err`
	}{
		{"123456", 123456, true},
		{"-123456", -123456, true},
		{"0", 0, true},
		{"1000000000", 1000000000, true},
		{"invalid", 0, false}, // Теперь проверяется `false`, если ошибка парсинга
	}

	for _, test := range tests {
		t.Run("Test ParseInt64", func(t *testing.T) {
			result, ok := ParseInt64(test.value)
			if test.ok {
				assert.True(t, ok, "Expected successful parse")
				assert.Equal(t, test.expected, result, "they should be equal")
			} else {
				assert.False(t, ok, "Expected parse to fail")
				assert.Equal(t, test.expected, result, "they should be equal")
			}
		})
	}
}
