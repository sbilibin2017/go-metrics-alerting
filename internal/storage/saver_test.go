package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaver_Save(t *testing.T) {
	// Create a new storage instance
	storage := NewStorage[string]()

	// Create a new Saver for the storage
	saver := NewSaver(storage)

	// Test data to save
	key := "testKey"
	value := "testValue"

	// Save data
	saver.Save(key, value)

	// Assert that the value has been saved correctly
	assert.Equal(t, value, storage.data[key], "The saved value should match the expected value")
}
