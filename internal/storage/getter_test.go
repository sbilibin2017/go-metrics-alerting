package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetter_Get(t *testing.T) {
	// Create a new storage instance
	storage := NewStorage[string]()

	// Create a new Saver to save data
	saver := NewSaver(storage)

	// Create a new Getter to retrieve data
	getter := NewGetter(storage)

	// Test data
	key := "testKey"
	value := "testValue"

	// Save data using Saver
	saver.Save(key, value)

	// Get data using Getter
	retrievedValue, exists := getter.Get(key)

	// Assert that the value exists and is correct
	assert.True(t, exists, "The value should exist")
	assert.Equal(t, value, retrievedValue, "The retrieved value should match the expected value")

	// Test non-existent key
	_, exists = getter.Get("nonExistentKey")
	assert.False(t, exists, "The non-existent key should return false")
}
