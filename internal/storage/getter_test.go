package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестирование функции NewGetter
func TestNewGetter(t *testing.T) {
	storage := NewStorage()
	getter := NewGetter(storage)

	// Проверяем, что getter не nil
	require.NotNil(t, getter, "Expected non-nil getter")

	// Проверяем, что в getter есть доступ к хранилищу
	assert.NotNil(t, getter.storage, "Expected getter to have a storage")
}

// Тестирование метода Get с существующим ключом
func TestGetter_Get_ExistingKey(t *testing.T) {
	storage := NewStorage()
	storage.data["key"] = "value"
	getter := NewGetter(storage)

	// Получаем значение по существующему ключу
	value, exists := getter.Get("key")

	// Проверяем, что значение и флаг существуют
	assert.True(t, exists, "Expected key 'key' to exist")
	assert.Equal(t, "value", value, "Expected 'value' for key 'key'")
}

// Тестирование метода Get с несуществующим ключом
func TestGetter_Get_NonExistingKey(t *testing.T) {
	storage := NewStorage()
	getter := NewGetter(storage)

	// Получаем значение по несуществующему ключу
	value, exists := getter.Get("nonexistent_key")

	// Проверяем, что ключ не существует
	assert.False(t, exists, "Expected key 'nonexistent_key' not to exist")
	assert.Empty(t, value, "Expected empty string for nonexistent key")
}

// Тестирование метода Get после добавления данных
func TestGetter_Get_AfterSave(t *testing.T) {
	storage := NewStorage()
	getter := NewGetter(storage)

	// Сохраняем пару ключ-значение
	storage.data["key"] = "value"

	// Получаем значение после добавления
	value, exists := getter.Get("key")

	// Проверяем, что данные успешно получены
	assert.True(t, exists, "Expected key 'key' to exist")
	assert.Equal(t, "value", value, "Expected 'value' for key 'key'")
}

// Тестирование метода Get с пустым значением
func TestGetter_Get_EmptyValue(t *testing.T) {
	storage := NewStorage()
	getter := NewGetter(storage)

	// Сохраняем пару ключ-значение с пустым значением
	storage.data["key"] = ""

	// Получаем значение
	value, exists := getter.Get("key")

	// Проверяем, что ключ существует, но значение пустое
	assert.True(t, exists, "Expected key 'key' to exist")
	assert.Equal(t, "", value, "Expected empty string for key 'key'")
}
