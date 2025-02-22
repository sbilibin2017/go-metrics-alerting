package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестирование функции NewRanger
func TestNewRanger(t *testing.T) {
	storage := NewStorage()
	ranger := NewRanger(storage)

	// Проверяем, что ranger не nil
	require.NotNil(t, ranger, "Expected non-nil ranger")

	// Проверяем, что в ranger есть доступ к хранилищу
	assert.NotNil(t, ranger.storage, "Expected ranger to have a storage")
}

// Тестирование метода Range с несколькими элементами
func TestRanger_Range_MultipleElements(t *testing.T) {
	storage := NewStorage()
	ranger := NewRanger(storage)

	// Добавляем несколько элементов в хранилище
	storage.data["key1"] = "value1"
	storage.data["key2"] = "value2"
	storage.data["key3"] = "value3"

	// Подсчитываем количество вызовов колбэка
	count := 0
	ranger.Range(func(key, value string) bool {
		count++
		assert.Contains(t, []string{"key1", "key2", "key3"}, key, "Expected key to be in storage")
		assert.Contains(t, []string{"value1", "value2", "value3"}, value, "Expected value to be in storage")
		return true
	})

	// Проверяем, что колбэк был вызван для всех элементов
	assert.Equal(t, 3, count, "Expected callback to be called for 3 elements")
}

// Тестирование метода Range с досрочным завершением (когда колбэк возвращает false)
func TestRanger_Range_EarlyExit(t *testing.T) {
	storage := NewStorage()
	ranger := NewRanger(storage)

	// Добавляем несколько элементов в хранилище
	storage.data["key1"] = "value1"
	storage.data["key2"] = "value2"
	storage.data["key3"] = "value3"

	// Подсчитываем количество вызовов колбэка
	count := 0
	ranger.Range(func(key, value string) bool {
		count++
		// Прерываем перебор после первого элемента
		return false
	})

	// Проверяем, что колбэк был вызван только один раз
	assert.Equal(t, 1, count, "Expected callback to be called only once")
}

// Тестирование метода Range с пустым хранилищем
func TestRanger_Range_EmptyStorage(t *testing.T) {
	storage := NewStorage()
	ranger := NewRanger(storage)

	// Перебор в пустом хранилище
	count := 0
	ranger.Range(func(key, value string) bool {
		count++
		return true
	})

	// Проверяем, что колбэк не был вызван
	assert.Equal(t, 0, count, "Expected callback not to be called for empty storage")
}

// Тестирование метода Range с одним элементом
func TestRanger_Range_SingleElement(t *testing.T) {
	storage := NewStorage()
	ranger := NewRanger(storage)

	// Добавляем один элемент в хранилище
	storage.data["key1"] = "value1"

	// Перебор с одним элементом
	count := 0
	ranger.Range(func(key, value string) bool {
		count++
		assert.Equal(t, "key1", key, "Expected key 'key1'")
		assert.Equal(t, "value1", value, "Expected value 'value1'")
		return true
	})

	// Проверяем, что колбэк был вызван один раз
	assert.Equal(t, 1, count, "Expected callback to be called once")
}
