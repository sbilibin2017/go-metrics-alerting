package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Set(t *testing.T) {
	s := NewMemStorage()
	s.Set("key1", "value1")

	val, ok := s.Get("key1")
	assert.True(t, ok, "expected key1 to exist")
	assert.Equal(t, "value1", val, "expected value1")
}

func TestStorage_Get(t *testing.T) {
	s := NewMemStorage()
	s.Set("key1", "value1")

	val, ok := s.Get("key1")
	assert.True(t, ok, "expected key1 to exist")
	assert.Equal(t, "value1", val, "expected value1")

	_, ok = s.Get("key2")
	assert.False(t, ok, "expected key2 to not exist")
}

func TestStorage_Range(t *testing.T) {
	s := NewMemStorage()
	s.Set("key1", "value1")
	s.Set("key2", "value2")

	keys := make(map[string]bool)
	s.Range(func(key, value string) bool {
		keys[key] = true
		return true
	})

	assert.Len(t, keys, 2, "expected two keys")
	assert.Contains(t, keys, "key1", "expected key1 in range")
	assert.Contains(t, keys, "key2", "expected key2 in range")
}

func TestStorage_RangeBreak(t *testing.T) {
	s := NewMemStorage()
	s.Set("key1", "value1")
	s.Set("key2", "value2")
	s.Set("key3", "value3")

	keys := make([]string, 0)

	// Break the iteration after the first key
	s.Range(func(key, value string) bool {
		keys = append(keys, key)
		return len(keys) < 2 // Stop after 2 iterations
	})

	assert.Len(t, keys, 2, "expected two keys")
	assert.Contains(t, keys, "key1", "expected key1 in range")
	assert.Contains(t, keys, "key2", "expected key2 in range")
}
