package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage_StringKeyIntValue(t *testing.T) {
	// Создаем хранилище string -> int
	store := NewStorage[string, int]()

	// Создание Saver и сохранение данных
	saver := NewSaver(store)
	saver.Save("key1", 100)
	saver.Save("key2", 200)

	// Создание Getter и проверка значений
	getter := NewGetter(store)

	// Тестируем существующие ключи
	value, exists := getter.Get("key1")
	require.True(t, exists)
	assert.Equal(t, 100, value)

	value, exists = getter.Get("key2")
	require.True(t, exists)
	assert.Equal(t, 200, value)

	// Тестируем несуществующий ключ
	value, exists = getter.Get("key3")
	assert.False(t, exists)
	assert.Equal(t, 0, value) // Проверяем, что возвращается нулевое значение для типа int
}

func TestStorage_Range(t *testing.T) {
	// Создаем хранилище string -> int
	store := NewStorage[string, int]()

	// Создание Saver и сохранение данных
	saver := NewSaver(store)
	saver.Save("key1", 100)
	saver.Save("key2", 200)
	saver.Save("key3", 300)

	// Создаем Ranger для перебора
	ranger := NewRanger(store)

	// Проверка перебора всех элементов
	expected := map[string]int{
		"key1": 100,
		"key2": 200,
		"key3": 300,
	}

	ranger.Range(func(key string, value int) bool {
		expectedValue, ok := expected[key]
		require.True(t, ok) // Ожидаем, что ключ будет найден в ожидаемой мапе
		assert.Equal(t, expectedValue, value)
		return true
	})
}

func TestStorage_EmptyStorage(t *testing.T) {
	// Создаем пустое хранилище
	store := NewStorage[string, int]()

	// Проверка, что хранилище пустое
	getter := NewGetter(store)
	_, exists := getter.Get("key1")
	assert.False(t, exists)
}

func TestStorage_DifferentTypes(t *testing.T) {
	// Создаем хранилище string -> float64
	store := NewStorage[string, float64]()

	// Создание Saver и сохранение данных
	saver := NewSaver(store)
	saver.Save("pi", 3.14159)
	saver.Save("e", 2.71828)

	// Создание Getter и проверка значений
	getter := NewGetter(store)

	// Проверяем сохраненные данные
	value, exists := getter.Get("pi")
	require.True(t, exists)
	assert.Equal(t, 3.14159, value)

	value, exists = getter.Get("e")
	require.True(t, exists)
	assert.Equal(t, 2.71828, value)
}

func TestStorage_NilValue(t *testing.T) {
	// Создаем хранилище string -> *string (указатель на строку)
	store := NewStorage[string, *string]()

	// Создаем Saver и сохраняем nil
	saver := NewSaver(store)
	var str *string
	saver.Save("key1", str)

	// Создание Getter и проверка значения
	getter := NewGetter(store)
	value, exists := getter.Get("key1")
	require.True(t, exists)
	assert.Nil(t, value) // Проверяем, что возвращается nil
}

func TestStorage_RangeWithCallbackExit(t *testing.T) {
	// Создаем хранилище string -> int
	store := NewStorage[string, int]()

	// Создание Saver и сохранение данных
	saver := NewSaver(store)
	saver.Save("key1", 100)
	saver.Save("key2", 200)
	saver.Save("key3", 300)

	// Создаем Ranger для перебора
	ranger := NewRanger(store)

	// Флаг для проверки, что цикл остановился
	callbackCalled := false

	// Выполняем перебор с условием остановки в callback
	ranger.Range(func(key string, value int) bool {
		// Проверяем первый элемент
		if key == "key2" {
			callbackCalled = true
			return false // Останавливаем перебор
		}
		return true // Продолжаем перебор
	})

	// Проверяем, что callback был вызван для "key2" и цикл остановился
	assert.True(t, callbackCalled, "Callback for key2 should have been called and loop should stop")
}
