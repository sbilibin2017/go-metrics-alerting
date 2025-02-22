package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестирование функции NewStorage
func TestNewStorage(t *testing.T) {
	storage := NewStorage()

	// Проверяем, что хранилище не nil
	require.NotNil(t, storage, "Expected non-nil storage")

	// Проверяем, что начальная карта данных пустая
	assert.Empty(t, storage.data, "Expected empty data map")
}

// Тестирование метода, который добавляет элемент в хранилище
func TestStorage_Add(t *testing.T) {
	storage := NewStorage()

	// Добавляем пару ключ-значение
	storage.data["key"] = "value"

	// Проверяем, что значение по ключу "key" равно "value"
	assert.Equal(t, "value", storage.data["key"], "Expected value 'value' for key 'key'")
}

// Тестирование метода чтения данных с блокировкой
func TestStorage_Read(t *testing.T) {
	storage := NewStorage()

	// Добавляем пару ключ-значение
	storage.data["key"] = "value"

	// Читаем значение
	storage.mu.RLock() // читающая блокировка
	defer storage.mu.RUnlock()
	val, ok := storage.data["key"]

	// Проверяем, что значение корректное
	assert.True(t, ok, "Expected key 'key' to exist")
	assert.Equal(t, "value", val, "Expected value 'value' for key 'key'")
}

// Тестирование метода записи с блокировкой
func TestStorage_Write(t *testing.T) {
	storage := NewStorage()

	// Блокируем запись
	storage.mu.Lock()
	storage.data["key"] = "new value"
	storage.mu.Unlock()

	// Проверяем, что значение было обновлено
	val, ok := storage.data["key"]
	assert.True(t, ok, "Expected key 'key' to exist")
	assert.Equal(t, "new value", val, "Expected value 'new value' for key 'key'")
}
