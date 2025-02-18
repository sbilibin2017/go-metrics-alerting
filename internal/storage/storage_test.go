package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)
	success := saver.Save("key1", "value1")
	assert.True(t, success, "Save operation should return true")
}

func TestGet(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)
	getter := NewGetter(storage)

	saver.Save("key1", "value1")

	value, exists := getter.Get("key1")
	assert.True(t, exists, "key1 should exist")
	assert.Equal(t, "value1", value, "Expected value1")

	_, exists = getter.Get("key2")
	assert.False(t, exists, "key2 should not exist")
}

func TestRange(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)
	ranger := NewRanger(storage)

	saver.Save("key1", "value1")
	saver.Save("key2", "value2")
	saver.Save("key3", "value3")

	found := make(map[string]bool)
	stoppedAt := "key2"

	ranger.Range(func(key, value string) bool {
		found[key] = true
		if key == stoppedAt {
			return false
		}
		return true
	})

	assert.True(t, found["key1"], "Expected to find key1")
	assert.True(t, found["key2"], "Expected to find key2")
	assert.False(t, found["key3"], "Expected to not find key3 as iteration should stop at key2")
}
