package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage_Creation(t *testing.T) {
	storage := NewStorage[string, int]()
	assert.NotNil(t, storage)
	assert.Empty(t, storage.data)
}
