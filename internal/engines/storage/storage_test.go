package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorageEngine_SetGet(t *testing.T) {
	se := NewStorageEngine[string, int]()

	se.Set("key1", 10)
	se.Set("key2", 20)

	val, ok := se.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 10, val)

	val, ok = se.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, 20, val)

	val, ok = se.Get("key3")
	assert.False(t, ok)
	assert.Equal(t, 0, val) // Ожидаемое значение по умолчанию для int
}

func TestStorageEngine_Generate(t *testing.T) {
	se := NewStorageEngine[string, int]()

	se.Set("key1", 10)
	se.Set("key2", 20)

	values := make(map[int]bool)
	for pair := range se.Generate() {
		values[pair[0]] = true
	}

	assert.True(t, values[10])
	assert.True(t, values[20])
	assert.Len(t, values, 2)
}
