package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	// Создаём новое хранилище для данных типа int
	storage := NewStorage[int]()

	// Проверяем, что хранилище не равно nil
	require.NotNil(t, storage)
	// Проверяем, что карта данных (map) инициализирована
	assert.NotNil(t, storage.data)
}
