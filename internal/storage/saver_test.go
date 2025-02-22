package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тестирование функции NewSaver
func TestNewSaver(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)

	// Проверяем, что saver не nil
	require.NotNil(t, saver, "Expected non-nil saver")

	// Проверяем, что в saver есть доступ к хранилищу
	assert.NotNil(t, saver.storage, "Expected saver to have a storage")
}

// Тестирование метода Save
func TestSaver_Save(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)

	// Сохраняем пару ключ-значение
	success := saver.Save("key", "value")

	// Проверяем, что операция завершена успешно
	assert.True(t, success, "Expected Save to return true")

	// Проверяем, что данные сохранены в хранилище
	assert.Equal(t, "value", storage.data["key"], "Expected 'value' for key 'key'")
}

// Тестирование метода Save с несколькими записями
func TestSaver_MultipleSaves(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)

	// Сохраняем несколько пар ключ-значение
	saver.Save("key1", "value1")
	saver.Save("key2", "value2")
	saver.Save("key3", "value3")

	// Проверяем, что все данные сохранены корректно
	assert.Equal(t, "value1", storage.data["key1"], "Expected 'value1' for key 'key1'")
	assert.Equal(t, "value2", storage.data["key2"], "Expected 'value2' for key 'key2'")
	assert.Equal(t, "value3", storage.data["key3"], "Expected 'value3' for key 'key3'")
}

// Тестирование метода Save при попытке сохранить пустое значение
func TestSaver_SaveEmptyValue(t *testing.T) {
	storage := NewStorage()
	saver := NewSaver(storage)

	// Сохраняем пару ключ-значение с пустым значением
	success := saver.Save("key", "")

	// Проверяем, что операция завершена успешно
	assert.True(t, success, "Expected Save to return true when saving empty value")

	// Проверяем, что пустое значение сохранено
	assert.Equal(t, "", storage.data["key"], "Expected empty string for key 'key'")
}
